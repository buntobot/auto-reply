package affinity

import (
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestTeamRandomCaptainLogins(t *testing.T) {
	team := Team{Captains: []*github.User{
		{Login: github.String("SuriyaaKudoIsc")},
		{Login: github.String("aahashderuffy")},
		{Login: github.String("subins2000")},
	}}
	selections := team.RandomCaptainLogins(1)
	assert.Len(t, selections, 1)
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[0])

	selections = team.RandomCaptainLogins(2)
	assert.Len(t, selections, 2)
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[0])
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[1])

	selections = team.RandomCaptainLogins(3)
	assert.Len(t, selections, 3)
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[0])
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[1])
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[2])

	selections = team.RandomCaptainLogins(4)
	assert.Len(t, selections, 3)
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[0])
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[1])
	assert.Contains(t, []string{"SuriyaaKudoIsc", "aahashderuffy", "subins2000"}, selections[2])
}

func TestTeamRandomCaptainLoginsExcluding(t *testing.T) {
	excluded := "SuriyaaKudoIsc"
	team := Team{Captains: []*github.User{
		{Login: github.String("SuriyaaKudoIsc")},
		{Login: github.String("aahashderuffy")},
		{Login: github.String("subins2000")},
	}}

	selections := team.RandomCaptainLoginsExcluding(excluded, 1)
	assert.Len(t, selections, 1)
	assert.Contains(t, []string{"aahashderuffy", "subins2000"}, selections[0])

	selections = team.RandomCaptainLoginsExcluding(excluded, 2)
	assert.Len(t, selections, 2)
	assert.Contains(t, []string{"aahashderuffy", "subins2000"}, selections[0])
	assert.Contains(t, []string{"aahashderuffy", "subins2000"}, selections[1])

	selections = team.RandomCaptainLoginsExcluding(excluded, 3)
	assert.Len(t, selections, 2)
	assert.Contains(t, []string{"aahashderuffy", "subins2000"}, selections[0])
	assert.Contains(t, []string{"aahashderuffy", "subins2000"}, selections[1])
}
