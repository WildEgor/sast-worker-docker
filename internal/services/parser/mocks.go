package parser

import (
	"github.com/stretchr/testify/mock"
	"os"
)

var _ IParser = (*MockDockerParserService)(nil)

// MockDockerParserService is a mock implementation of the DockerParserService.
type MockDockerParserService struct {
	mock.Mock
}

func (m *MockDockerParserService) Parse(file *os.File) (*DockerfileParsedResult, error) {
	args := m.Called(file)
	return args.Get(0).(*DockerfileParsedResult), args.Error(1)
}
