package client

import "github.com/stretchr/testify/mock"

type MockClient struct {
	mock.Mock
}

var _ Client = (*MockClient)(nil)

func (m *MockClient) CurrentBranch() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockClient) RefLog(cfg RefLogConfig) ([]string, error) {
	args := m.Called(cfg)
	return args.Get(0).([]string), args.Error(1)
}

func (m *MockClient) Checkout(b string) error {
	return m.Called(b).Error(0)
}
