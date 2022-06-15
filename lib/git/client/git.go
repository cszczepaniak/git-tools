package client

type Client interface {
	CurrentBranch() (string, error)
	RefLog(cfg RefLogConfig) ([]string, error)
	Checkout(b string) error
}

type RefLogConfig struct {
	Pretty string
	Date   string
	Count  int
}
