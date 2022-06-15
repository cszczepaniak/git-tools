package git

import (
	"fmt"
	"strings"

	"github.com/cszczepaniak/git-tools/lib/git/client"
)

type LatestBranch struct {
	name      string
	timestamp string
}

func (b LatestBranch) Name() string {
	return b.name
}

func (b LatestBranch) String() string {
	return fmt.Sprintf(`%s (%s)`, b.name, b.timestamp)
}

func LatestBranches(c client.Client, refLogLimit int) ([]LatestBranch, error) {
	const (
		toDelim         = ` to `
		refLogPartDelim = ` ~ `
		checkoutPrefix  = `checkout: `
		leftCurly       = `{`
		rightCurly      = `}`
	)

	refLog, err := c.RefLog(client.RefLogConfig{
		Pretty: `%gs` + refLogPartDelim + `%gd`,
		Date:   `relative`,
		Count:  refLogLimit,
	})
	if err != nil {
		return nil, err
	}

	current, err := c.CurrentBranch()
	if err != nil {
		return nil, err
	}

	branchSet := make(map[string]struct{})
	var latest []LatestBranch
	for _, entry := range refLog {
		if !strings.HasPrefix(entry, checkoutPrefix) {
			continue
		}

		parts := strings.SplitN(entry[len(checkoutPrefix):], refLogPartDelim, 2)
		if len(parts) != 2 {
			continue
		}

		branchMove := parts[0]
		to := strings.Index(branchMove, toDelim)
		if to < 0 {
			continue
		}

		branch := branchMove[to+len(toDelim):]
		switch branch {
		case `master`, `main`, current:
			continue
		}

		if _, ok := branchSet[branch]; ok {
			continue
		}

		timestampStr := parts[1] // HEAD@{timestamp}
		leftCurlyIdx := strings.Index(timestampStr, leftCurly)
		rightCurlyIdx := strings.Index(timestampStr, rightCurly)

		if leftCurlyIdx < 0 || rightCurlyIdx < 0 || rightCurlyIdx <= leftCurlyIdx {
			continue
		}

		branchSet[branch] = struct{}{}
		latest = append(latest, LatestBranch{
			name:      branch,
			timestamp: timestampStr[leftCurlyIdx+1 : rightCurlyIdx],
		})
	}

	return latest, nil
}
