package client

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupIntegrationTest(t *testing.T) string {
	dir, err := os.MkdirTemp(``, ``)
	require.NoError(t, err)

	runCmd(t, dir, `git`, `init`)
	runCmd(t, dir, `git`, `commit`, `--allow-empty`, `-m`, `initial commit`)

	t.Cleanup(func() {
		err := os.RemoveAll(dir)
		if err != nil {
			t.Error(`error removing temp dir: `, dir)
		}
	})

	return dir
}

func TestCurrentBranch(t *testing.T) {
	dir := setupIntegrationTest(t)

	runCmd(t, dir, `git`, `checkout`, `-b`, `foo`)

	c := NewClient(dir)

	b, err := c.CurrentBranch()
	require.NoError(t, err)
	assert.Equal(t, `foo`, b)
}

func TestRefLog(t *testing.T) {
	dir := setupIntegrationTest(t)

	for i := 0; i < 3; i++ {
		runCmd(t, dir, `git`, `checkout`, `-b`, fmt.Sprintf(`branch%d`, i))
	}
	for i := 0; i < 3; i++ {
		runCmd(t, dir, `git`, `commit`, `--allow-empty`, `-m`, fmt.Sprintf(`foo%d`, i))
	}

	c := NewClient(dir)

	refLog, err := c.RefLog(RefLogConfig{
		Pretty: `%gd, %gs`,
		Date:   `relative`,
		Count:  5,
	})
	require.NoError(t, err)
	assert.Len(t, refLog, 5)

	expActions := []string{
		`commit: foo2`,
		`commit: foo1`,
		`commit: foo0`,
		`checkout: moving from branch1 to branch2`,
		`checkout: moving from branch0 to branch1`,
	}
	for i, entry := range refLog {
		parts := strings.SplitN(entry, `, `, 2)
		require.Len(t, parts, 2)
		assert.Regexp(t, `HEAD@\{\d+ seconds? ago\}`, parts[0])
		assert.Equal(t, expActions[i], parts[1])
	}
}

func TestCheckout(t *testing.T) {
	dir := setupIntegrationTest(t)

	runCmd(t, dir, `git`, `checkout`, `-b`, `branch0`)
	runCmd(t, dir, `git`, `checkout`, `-b`, `branch1`)
	runCmd(t, dir, `git`, `checkout`, `-b`, `branch2`)

	c := NewClient(dir)

	require.NoError(t, c.Checkout(`branch1`))
	current, err := c.CurrentBranch()
	require.NoError(t, err)
	assert.Equal(t, `branch1`, current)

	require.NoError(t, c.Checkout(`branch0`))
	current, err = c.CurrentBranch()
	require.NoError(t, err)
	assert.Equal(t, `branch0`, current)

	require.Error(t, c.Checkout(`dne`))
	current, err = c.CurrentBranch()
	require.NoError(t, err)
	assert.Equal(t, `branch0`, current)
}

func runCmd(t *testing.T, dir, name string, args ...string) string {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	require.NoError(t, err)

	return strings.TrimSpace(string(out))
}
