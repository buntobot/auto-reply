package affinity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleLongComment = `On the site documentation section, links to documentation sections always point to the buntorb.com website, this means that users testing changes might get confused because they will see the official external website page instead of their local website upon clicking those links.


**Please check if this change doesn't break the official website on https://buntorb.com before accepting the pull request.**

----------

@bunto/documentation`

func TestFindAffinityTeam(t *testing.T) {
	allTeams := []Team{
		{ID: 456, Mention: "@bunto/documentation"},
		{ID: 789, Mention: "@bunto/ecosystem"},
		{ID: 101, Mention: "@bunto/performance"},
		{ID: 213, Mention: "@bunto/stability"},
		{ID: 141, Mention: "@bunto/windows"},
		{ID: 123, Mention: "@bunto/build"},
	}

	examples := []struct {
		body           string
		matchingTeamID int
	}{
		{exampleLongComment, 456},
		{"@bunto/documentation @bunto/build", 456},
		{"@bunto/windows @bunto/documentation", 456},
		{"@bunto/windows", 141},
	}
	for _, example := range examples {
		matchingTeam, err := findAffinityTeam(example.body, allTeams)
		assert.NoError(t, err)
		assert.Equal(t, matchingTeam.ID, example.matchingTeamID,
			"expected the following to match %d team: `%s`", example.matchingTeamID, example.body)
	}
}
