package git

import (
	"strings"

	"github.com/cszczepaniak/git-tools/lib/git/client"
)

func LatestBranches(c client.Client, refLogLimit int) (map[string]string, error) {
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

	branchToTime := make(map[string]string)
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

		if _, ok := branchToTime[branch]; ok {
			continue
		}

		timestampStr := parts[1] // HEAD@{timestamp}
		leftCurlyIdx := strings.Index(timestampStr, leftCurly)
		rightCurlyIdx := strings.Index(timestampStr, rightCurly)

		if leftCurlyIdx < 0 || rightCurlyIdx < 0 || rightCurlyIdx <= leftCurlyIdx {
			continue
		}

		branchToTime[branch] = timestampStr[leftCurlyIdx+1 : rightCurlyIdx]
	}

	return branchToTime, nil
}
