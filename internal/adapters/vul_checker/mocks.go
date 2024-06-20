package vul_checker

import (
	"github.com/stretchr/testify/mock"
)

var _ IVulChecker = (*MockVulChecker)(nil)

// MockVulChecker is a mock implementation of the ILinter interface.
type MockVulChecker struct {
	mock.Mock
}

func (m *MockVulChecker) GetCVEByImageNames(names []string) []VulListItem {
	args := m.Called(names)
	return args.Get(0).([]VulListItem)
}
