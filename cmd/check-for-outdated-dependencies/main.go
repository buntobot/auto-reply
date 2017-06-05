// check-for-outdated-dependencies takes a repo
package main

import (
	"flag"
	"log"
	"strings"

	"github.com/buntobot/auto-reply/ctx"
	"github.com/buntobot/auto-reply/dependencies"
)

var defaultRepos = strings.Join([]string{
	"bunto/bunto",
}, ",")

func main() {
	var depType string
	flag.StringVar(&depType, "type", "ruby", "The type of dependency we're checking (options: ruby)")
	var reposString string
	flag.StringVar(&reposString, "repos", defaultRepos, "Comma-separated list of repos to check, e.g. bunto/bunto,bunto/bunto-import")
	var perform bool
	flag.BoolVar(&perform, "f", false, "Whether to open issues (default: false, which is a dry-run)")
	flag.Parse()

	context := ctx.NewDefaultContext()

	for _, repo := range strings.Split(reposString, ",") {
		pieces := strings.SplitN(repo, "/", 2)
		repoOwner, repoName := pieces[0], pieces[1]
		checker := dependencies.NewRubyDependencyChecker(repoOwner, repoName)
		outdated := checker.AllOutdatedDependencies(context)
		for _, dependency := range outdated {
			log.Printf(
				"%s/%s: %s is outdated (constraint: %s, but latest version is %s)",
				repoOwner, repoName, dependency.GetName(), dependency.GetConstraint(), dependency.GetLatestVersion(context),
			)

			// Do not open issues if dry-run.
			if !perform {
				continue
			}

			preExistingIssue := dependencies.GitHubUpdateIssueForDependency(context, repoOwner, repoName, dependency)

			if preExistingIssue == nil {
				issue, err := dependencies.FileGitHubIssueForDependency(context, repoOwner, repoName, dependency)
				if err != nil {
					log.Printf("%s/%s: error creating issue for %s: %v", repoOwner, repoName, dependency.GetName(), err)
				} else {
					log.Printf("%s/%s: issue for %s filed: %s", repoOwner, repoName, dependency.GetName(), *issue.HTMLURL)
				}
			} else {
				log.Printf("%s/%s: issue for %s already open: %s",
					repoOwner, repoName, dependency.GetName(), *preExistingIssue.HTMLURL)
			}
		}
	}
}
