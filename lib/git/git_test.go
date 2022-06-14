package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/cszczepaniak/git-tools/lib/git/client"
)

func TestLatestBranches(t *testing.T) {
	c := &client.MockClient{}

	c.On(`RefLog`, mock.Anything).Return([]string{
		`commit: here's a commit ~ HEAD@{1 days ago}`,
		`checkout: moving from branch1 to branch2 ~ HEAD@{2 days ago}`,
		`checkout: moving from branch1 too branch2 ~ HEAD@{3 days ago}`, // ignore malformed entry
		`checkout: moving from branch1 to branch2 -- HEAD@{3 days ago}`, // ignore malformed entry
		`checkout: moving from branch1 to branch2 ~ HEAD@}3 days ago{`,  // ignore malformed entry
		`checkout: moving from branch1 to branch2 ~ HEAD@{3 days ago`,   // ignore malformed entry
		`checkout: moving from branch1 to branch2 ~ HEAD{3 days ago}`,   // ignore malformed entry
		`checkout: moving from branch1 to branch2 ~ HEAD@{4 days ago}`,  // duplicates should not be included
		`checkout: moving from branch1 to main ~ HEAD@{5 days ago}`,     // ignore main branch
		`checkout: moving from branch1 to master ~ HEAD@{6 days ago}`,   // ignore main branch
		`something: blah ~ HEAD@{7 days ago}`,
		`checkout: moving from branch2 to branch1 ~ HEAD@{8 days ago}`,
		`checkout: moving from branch1 to current ~ HEAD@{9 days ago}`, // current ignored
	}, nil).Once()
	c.On(`CurrentBranch`).Return(`current`, nil).Once()

	latest, err := LatestBranches(c, 1000)
	require.NoError(t, err)

	assert.Equal(t, map[string]string{
		`branch2`: `2 days ago`,
		`branch1`: `8 days ago`,
	}, latest)

	c.AssertExpectations(t)

	c = &client.MockClient{}
	c.On(`RefLog`, mock.MatchedBy(func(cfg client.RefLogConfig) bool {
		return cfg.Count == 1234
	})).Return([]string{}, nil).Once()
	c.On(`CurrentBranch`).Return(`blah`, nil).Once()

	latest, err = LatestBranches(c, 1234)
	require.NoError(t, err)
	require.Empty(t, latest)

	c.AssertExpectations(t)
}
