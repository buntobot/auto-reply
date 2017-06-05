package lgtm

import (
	"fmt"
	"testing"

	"github.com/google/go-github/github"
	"github.com/stretchr/testify/assert"
)

func TestParseStatus(t *testing.T) {
	cases := []struct {
		sha             string
		description     string
		expectedLgtmers []string
		expectedQuorum  int
	}{
		{"deadbeef", "", []string{}, 0},
		{"deadbeef", "Waiting for approval from at least 2 maintainers.", []string{}, 2},
		// {"deadbeef", "Waiting for approval from at least 22 maintainers.", []string{}, 22},
		{"deadbeef", "Approved by @SuriyaaKudoIsc. Requires 1 more LGTM.", []string{"@SuriyaaKudoIsc"}, 2},
		{"deadbeef", "@SuriyaaKudoIsc have approved this PR. Requires 32 more LGTM's.", []string{"@SuriyaaKudoIsc"}, 33},
		// {"deadbeef", "@SuriyaaKudoIsc, and @aahashderuffy have approved this PR.", []string{"@SuriyaaKudoIsc", "@aahashderuffy"}, 2},
		// {"deadbeef", "@subins2000, @SuriyaaKudoIsc have approved this PR. Requires no more LGTM's.", []string{"@subins2000", "@SuriyaaKudoIsc", "@aahashderuffy"}, 3},
	}
	for _, test := range cases {
		parsed := parseStatus(test.sha, &github.RepoStatus{Description: github.String(test.description)})
		assert.Equal(t,
			test.expectedLgtmers, parsed.lgtmers,
			fmt.Sprintf("parsing description: %q", test.description))
		assert.Equal(t,
			test.expectedQuorum, parsed.quorum,
			fmt.Sprintf("parsing description: %q", test.description))
		assert.Equal(t, test.sha, parsed.sha)
	}
}

func TestStatusInfoIsLGTMer(t *testing.T) {
	cases := []struct {
		info             statusInfo
		lgtmerInQuestion string
		islgtmer         bool
	}{
		{statusInfo{}, "@SuriyaaKudoIsc", false},
		{statusInfo{lgtmers: []string{"@SuriyaaKudoIsc"}}, "@SuriyaaKudoIsc", true},
		// {statusInfo{lgtmers: []string{"@SuriyaaKudoIsc"}}, "@subins2000", false},
		// {statusInfo{lgtmers: []string{"@SuriyaaKudoIsc", "@subins2000"}}, "@subins2000", true},
		// {statusInfo{lgtmers: []string{"@SuriyaaKudoIsc", "@subins2000"}}, "@SuriyaaKudoIsc-", false},
		// {statusInfo{lgtmers: []string{"@SuriyaaKudoIsc", "@subins2000"}}, "@SuriyaaKudoIsc", true},
		// {statusInfo{lgtmers: []string{"@SuriyaaKudoIsc", "@subins2000"}}, "@SURIYAAKUDOISC", true},
		// {statusInfo{lgtmers: []string{"@aahashderuffy", "@subins2000"}}, "@aahashderuffy", true},
	}
	for _, test := range cases {
		assert.Equal(t,
			test.islgtmer, test.info.IsLGTMer(test.lgtmerInQuestion),
			fmt.Sprintf("asking about: %q for lgtmers: %q", test.lgtmerInQuestion, test.info.lgtmers))
	}
}

func TestNewState(t *testing.T) {
	cases := []struct {
		lgtmers  []string
		quorum   int
		expected string
	}{
		{[]string{}, 0, "success"},
		{[]string{}, 1, "pending"},
		{[]string{}, 2, "pending"},
		{[]string{"@SuriyaaKudoIsc"}, 0, "success"},
		{[]string{"@SuriyaaKudoIsc"}, 1, "success"},
		{[]string{"@SuriyaaKudoIsc"}, 2, "pending"},
		// {[]string{"@SuriyaaKudoIsc", "@subins2000"}, 0, "success"},
		// {[]string{"@SuriyaaKudoIsc", "@subins2000"}, 1, "success"},
		// {[]string{"@SuriyaaKudoIsc", "@subins2000"}, 2, "success"},
	}
	for _, test := range cases {
		info := statusInfo{lgtmers: test.lgtmers, quorum: test.quorum}
		assert.Equal(t,
			test.expected, info.newState(),
			fmt.Sprintf("with lgtmers: %q and quorum: %d", test.lgtmers, test.quorum))
	}
}

func TestNewDescription(t *testing.T) {
	cases := []struct {
		lgtmers     []string
		quorum      int
		description string
	}{
		{nil, 0, "No approval is required."},
		{nil, 1, "Awaiting approval from at least 1 maintainer."},
		// {[]string{}, 2, "Awaiting approval from at least 2 maintainers."},
		{[]string{"@SuriyaaKudoIsc"}, 2, "Approved by @SuriyaaKudoIsc. Requires 1 more LGTM."},
		// {[]string{"@SuriyaaKudoIsc", "@aahashderuffy"}, 2, "Approved by @SuriyaaKudoIsc and @aahashderuffy."},
		// {[]string{"@subins2000", "@aahashderuffy", "@SuriyaaKudoIsc"}, 5, "Approved by @subins2000, @aahashderuffy, and @SuriyaaKudoIsc. Requires 2 more LGTM's."},
	}
	for _, test := range cases {
		info := statusInfo{lgtmers: test.lgtmers, quorum: test.quorum}
		actual := info.newDescription()
		assert.Equal(t, test.description, actual)
		assert.True(t, len(actual) <= 140, fmt.Sprintf("%q must be <= 140 chars.", actual))
	}
}

func TestLGTMsRequiredDescription(t *testing.T) {
	cases := []struct {
		lgtmers  []string
		quorum   int
		expected string
	}{
		{nil, 0, ""},
		{nil, 1, "Requires 1 more LGTM."},
		{[]string{}, 2, "Requires 2 more LGTM's."},
		{[]string{"@SuriyaaKudoIsc"}, 2, "Requires 1 more LGTM."},
		// {[]string{"@SuriyaaKudoIsc", "@aahashderuffy"}, 2, ""},
		// {[]string{"@subins2000", "@aahashderuffy", "@SuriyaaKudoIsc"}, 5, "Requires 2 more LGTM's."},
	}
	for _, test := range cases {
		info := statusInfo{lgtmers: test.lgtmers, quorum: test.quorum}
		actual := info.newLGTMsRequiredDescription()
		assert.Equal(t, test.expected, actual)
		assert.True(t, len(actual) <= 140, fmt.Sprintf("%q must be <= 140 chars.", actual))
	}
}

func TestNewApprovedByDescription(t *testing.T) {
}

func TestStatusInfoNewRepoStatus(t *testing.T) {
	cases := []struct {
		owner          string
		lgtmers        []string
		quorum         int
		expContext     string
		expState       string
		expDescription string
	}{
		{"octocat", []string{}, 0, "octocat/lgtm", "success", "No approval is required."},
		{"SuriyaaKudoIsc", []string{}, 0, "SuriyaaKudoIsc/lgtm", "success", "No approval is required."},
		{"bunto", []string{}, 1, "bunto/lgtm", "pending", "Awaiting approval from at least 1 maintainer."},
		{"bunto", []string{"@SuriyaaKudoIsc"}, 1, "bunto/lgtm", "success", "Approved by @SuriyaaKudoIsc."},
		{"bunto", []string{"@SuriyaaKudoIsc"}, 2, "bunto/lgtm", "pending", "Approved by @SuriyaaKudoIsc. Requires 1 more LGTM."},
		// {"bunto", []string{"@SuriyaaKudoIsc", "@aahashderuffy"}, 1, "bunto/lgtm", "success", "Approved by @SuriyaaKudoIsc and @aahashderuffy."},
		// {"bunto", []string{"@SuriyaaKudoIsc", "@aahashderuffy"}, 2, "bunto/lgtm", "success", "Approved by @SuriyaaKudoIsc and @aahashderuffy."},
		// {"bunto", []string{"@SuriyaaKudoIsc", "@subins2000", "@aahashderuffy"}, 6, "bunto/lgtm", "pending", "Approved by @SuriyaaKudoIsc, @subins2000, and @aahashderuffy. Requires 3 more LGTM's."},
	}
	for _, test := range cases {
		status := statusInfo{lgtmers: test.lgtmers, quorum: test.quorum}
		newStatus := status.NewRepoStatus(test.owner)
		assert.Equal(t,
			test.expContext, *newStatus.Context,
			fmt.Sprintf("with lgtmers: %q and quorum: %d", test.lgtmers, test.quorum))
		assert.Equal(t,
			test.expState, *newStatus.State,
			fmt.Sprintf("with lgtmers: %q and quorum: %d", test.lgtmers, test.quorum))
		assert.Equal(t,
			test.expDescription, *newStatus.Description,
			fmt.Sprintf("with lgtmers: %q and quorum: %d", test.lgtmers, test.quorum))
		assert.True(t, len(*newStatus.Description) <= 140, fmt.Sprintf("%q must be <= 140 chars.", *newStatus.Description))
	}
}
