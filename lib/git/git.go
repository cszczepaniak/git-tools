package git

type Client interface {
	CurrentBranch() (string, error)
	RefLog(cfg RefLogConfig) ([]string, error)
}

type RefLogConfig struct {
	Pretty string
	Date   string
	Count  int
}
