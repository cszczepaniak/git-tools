package client

type Client interface {
	CurrentBranch() (string, error)
	Checkout(b string) error

	RefLog(cfg RefLogConfig) ([]string, error)

	ListConfigs(cfg ConfigConfig) (map[string]string, error)
}

type RefLogConfig struct {
	Pretty string
	Date   string
	Count  int
}

type ConfigConfig struct {
	Global bool
}
