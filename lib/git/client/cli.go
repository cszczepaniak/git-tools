package client

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func NewClient(dir string) Client {
	return &cliClient{
		dir: dir,
	}
}

type cliClient struct {
	dir string
}

func (c *cliClient) command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Dir = c.dir
	return cmd
}

func (c *cliClient) CurrentBranch() (string, error) {
	out, err := c.command(`git`, `rev-parse`, `--abbrev-ref`, `HEAD`).CombinedOutput()
	if err != nil {
		return ``, cliError(out, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func (cfg RefLogConfig) toArgs() []string {
	res := []string{`reflog`, `show`}
	if cfg.Pretty != `` {
		res = append(res, `--pretty=`+cfg.Pretty)
	}
	if cfg.Date != `` {
		res = append(res, `--date=`+cfg.Date)
	}
	if cfg.Count != 0 {
		res = append(res, `-n`, strconv.Itoa(cfg.Count))
	}
	return res
}

func (c *cliClient) RefLog(cfg RefLogConfig) ([]string, error) {
	cmd := c.command(`git`, cfg.toArgs()...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, cliError(stderr.Bytes(), err)
	}

	var res []string
	s := bufio.NewScanner(&stdout)
	for s.Scan() {
		res = append(res, s.Text())
	}
	return res, nil
}

func (c *cliClient) Checkout(b string) error {
	return cliError(c.command(`git`, `checkout`, b).CombinedOutput())
}

func cliError(out []byte, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w\n%s", err, out)
}
