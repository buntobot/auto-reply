// bunto is the configuration of handlers and such specific to the org's requirements. This is what you should copy and customize.
package bunto

import (
	"fmt"

	"github.com/buntobot/auto-reply/affinity"
	"github.com/buntobot/auto-reply/autopull"
	"github.com/buntobot/auto-reply/chlog"
	"github.com/buntobot/auto-reply/ctx"
	"github.com/buntobot/auto-reply/hooks"
	"github.com/buntobot/auto-reply/labeler"
	"github.com/buntobot/auto-reply/lgtm"
	"github.com/buntobot/auto-reply/travis"

	"github.com/google/go-github/github"
	"github.com/buntobot/auto-reply/bunto/deprecate"
	"github.com/buntobot/auto-reply/bunto/issuecomment"
)

var buntoOrgEventHandlers = hooks.EventHandlerMap{
	hooks.CreateEvent: {chlog.CreateReleaseOnTagHandler},
	hooks.IssuesEvent: {deprecate.DeprecateOldRepos},
	hooks.IssueCommentEvent: {
		issuecomment.PendingFeedbackUnlabeler,
		issuecomment.StaleUnlabeler,
		chlog.MergeAndLabel,
	},
	hooks.PullRequestEvent: {
		labeler.IssueHasPullRequestLabeler,
		labeler.PendingRebaseNeedsWorkPRUnlabeler,
	},
	hooks.StatusEvent: {statStatus, travis.FailingFmtBuildHandler},
}

func statStatus(context *ctx.Context, payload interface{}) error {
	status, ok := payload.(*github.StatusEvent)
	if !ok {
		return context.NewError("statStatus: not an status event")
	}

	context.SetIssue(*status.Repo.Owner.Login, *status.Repo.Name, -1)

	if context.Statsd != nil {
		statName := fmt.Sprintf("status.%s", *status.State)
		context.Log("context.Statsd.Count(%s, 1, []string{context:%s, repo:%s}, 1)", statName, *status.Context, context.Issue)
		return context.Statsd.Incr(
			statName,
			[]string{
				"context:" + *status.Context,
				"repo:" + context.Issue.String(),
			},
			float64(1.0), // rate..?
		)
	}
	return nil
}

func buntoAffinityHandler(context *ctx.Context) *affinity.Handler {
	handler := &affinity.Handler{}

	handler.AddRepo("bunto", "bunto")
	handler.AddRepo("bunto", "minima")

	handler.AddTeam(context, 1961060) // @bunto/build
	handler.AddTeam(context, 1961072) // @bunto/documentation
	handler.AddTeam(context, 1961061) // @bunto/ecosystem
	handler.AddTeam(context, 1961065) // @bunto/performance
	handler.AddTeam(context, 1961059) // @bunto/stability
	handler.AddTeam(context, 1116640) // @bunto/windows

	context.Log("affinity teams: %+v", handler.GetTeams())
	context.Log("affinity team repos: %+v", handler.GetRepos())

	return handler
}

func newLgtmHandler() *lgtm.Handler {
	handler := &lgtm.Handler{}

	handler.AddRepo("bunto", "bunto", 2)
	handler.AddRepo("bunto", "bunto-admin", 1)
	handler.AddRepo("bunto", "bunto-coffeescript", 2)
	handler.AddRepo("bunto", "bunto-compose", 1)
	handler.AddRepo("bunto", "bunto-docs", 1)
	handler.AddRepo("bunto", "bunto-feed", 1)
	handler.AddRepo("bunto", "bunto-gist", 2)
	handler.AddRepo("bunto", "bunto-import", 1)
	handler.AddRepo("bunto", "bunto-mentions", 2)
	handler.AddRepo("bunto", "bunto-opal", 2)
	handler.AddRepo("bunto", "bunto-paginate", 2)
	handler.AddRepo("bunto", "bunto-redirect-from", 2)
	handler.AddRepo("bunto", "bunto-sass-converter", 2)
	handler.AddRepo("bunto", "bunto-seo-tag", 1)
	handler.AddRepo("bunto", "bunto-sitemap", 2)
	handler.AddRepo("bunto", "bunto-textile-converter", 2)
	handler.AddRepo("bunto", "bunto-watch", 2)
	handler.AddRepo("bunto", "github-metadata", 2)
	handler.AddRepo("bunto", "jemoji", 1)
	handler.AddRepo("bunto", "mercenary", 1)
	handler.AddRepo("bunto", "minima", 1)
	handler.AddRepo("bunto", "plugins", 1)

	return handler
}

func NewBuntoOrgHandler(context *ctx.Context) *hooks.GlobalHandler {
	affinityHandler := buntoAffinityHandler(context)
	buntoOrgEventHandlers.AddHandler(hooks.IssuesEvent, affinityHandler.AssignIssueToAffinityTeamCaptain)
	buntoOrgEventHandlers.AddHandler(hooks.IssueCommentEvent, affinityHandler.AssignIssueToAffinityTeamCaptainFromComment)
	buntoOrgEventHandlers.AddHandler(hooks.PullRequestEvent, affinityHandler.AssignPRToAffinityTeamCaptain)

	lgtmHandler := newLgtmHandler()
	buntoOrgEventHandlers.AddHandler(hooks.PullRequestReviewEvent, lgtmHandler.PullRequestReviewHandler)

	autopullHandler := autopull.Handler{}
	autopullHandler.AcceptAllRepos(true)
	buntoOrgEventHandlers.AddHandler(hooks.PushEvent, autopullHandler.CreatePullRequestFromPush)

	return &hooks.GlobalHandler{
		Context:       context,
		EventHandlers: buntoOrgEventHandlers,
	}
}
