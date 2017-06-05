// auth provides a means of determining use permissions on GitHub.com for repositories.
package auth

import (
	"fmt"
	"log"

	"github.com/google/go-github/github"
	"github.com/buntobot/auto-reply/ctx"
)

var (
	teamsCache             = map[string][]*github.Team{}
	teamHasPushAccessCache = map[string]*github.Repository{}
	teamMembershipCache    = map[string]bool{}
	orgOwnersCache         = map[string][]*github.User{}
)

type authenticator struct {
	context *ctx.Context
}

func CommenterHasPushAccess(context *ctx.Context, event github.IssueCommentEvent) bool {
	auth := authenticator{context: context}
	orgTeams := auth.teamsForOrg(*event.Repo.Owner.Login)
	for _, team := range orgTeams {
		if auth.isTeamMember(*team.ID, *event.Comment.User.Login) &&
			auth.teamHasPushAccess(*team.ID, *event.Repo.Owner.Login, *event.Repo.Name) {
			return true
		}
	}
	return false
}

func UserIsOrgOwner(context *ctx.Context, org, login string) bool {
	auth := authenticator{context: context}
	for _, owner := range auth.ownersForOrg(org) {
		if *owner.Login == login {
			return true
		}
	}
	return false
}

func (auth authenticator) isTeamMember(teamId int, login string) bool {
	cacheKey := auth.cacheKeyIsTeamMember(teamId, login)
	if _, ok := teamMembershipCache[cacheKey]; !ok {
		newOk, _, err := auth.context.GitHub.Organizations.IsTeamMember(teamId, login)
		if err != nil {
			log.Printf("ERROR performing IsTeamMember(%d, \"%s\"): %v", teamId, login, err)
			return false
		}
		teamMembershipCache[cacheKey] = newOk
	}
	return teamMembershipCache[cacheKey]
}

func (auth authenticator) teamHasPushAccess(teamId int, owner, repo string) bool {
	cacheKey := auth.cacheKeyTeamHashPushAccess(teamId, owner, repo)
	if _, ok := teamHasPushAccessCache[cacheKey]; !ok {
		repository, _, err := auth.context.GitHub.Organizations.IsTeamRepo(teamId, owner, repo)
		if err != nil {
			log.Printf("ERROR performing IsTeamRepo(%d, \"%s\", \"%s\"): %v", teamId, owner, repo, err)
			return false
		}
		if repository == nil {
			return false
		}
		teamHasPushAccessCache[cacheKey] = repository
	}
	permissions := *teamHasPushAccessCache[cacheKey].Permissions
	return permissions["push"] || permissions["admin"]
}

func (auth authenticator) teamsForOrg(org string) []*github.Team {
	if _, ok := teamsCache[org]; !ok {
		teamz, _, err := auth.context.GitHub.Organizations.ListTeams(org, &github.ListOptions{
			Page: 0, PerPage: 100,
		})
		if err != nil {
			log.Printf("ERROR performing ListTeams(\"%s\"): %v", org, err)
			return nil
		}
		teamsCache[org] = teamz
	}
	return teamsCache[org]
}

func (auth authenticator) ownersForOrg(org string) []*github.User {
	if _, ok := orgOwnersCache[org]; !ok {
		owners, _, err := auth.context.GitHub.Organizations.ListMembers(org, &github.ListMembersOptions{
			Role: "admin", // owners
		})
		if err != nil {
			auth.context.Log("ERROR performing ListMembers(\"%s\"): %v", org, err)
			return nil
		}
		orgOwnersCache[org] = owners
	}
	return orgOwnersCache[org]
}

func (auth authenticator) cacheKeyIsTeamMember(teamId int, login string) string {
	return fmt.Sprintf("%d_%s", teamId, login)
}

func (auth authenticator) cacheKeyTeamHashPushAccess(teamId int, owner, repo string) string {
	return fmt.Sprintf("%d_%s_%s", teamId, owner, repo)
}
