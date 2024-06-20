package linter

import (
	"github.com/stretchr/testify/mock"
	"os"
)

var _ ILinter = (*MockLinter)(nil)

// MockLinter is a mock implementation of the ILinter interface.
type MockLinter struct {
	mock.Mock
}

func (m *MockLinter) Check(file *os.File) ([]CheckResult, error) {
	args := m.Called(file)
	return args.Get(0).([]CheckResult), args.Error(1)
}
