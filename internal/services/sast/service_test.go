package sast

import (
	"github.com/WildEgor/sast-worker-docker/internal/adapters/linter"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/vul_checker"
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"github.com/WildEgor/sast-worker-docker/internal/services/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestService_Analyze(t *testing.T) {
	appConfig := &configs.AppConfig{
		TempPath: os.TempDir(),
	}

	mockLinter := new(linter.MockLinter)
	mockChecker := new(vul_checker.MockVulChecker)
	mockDockerParser := new(parser.MockDockerParserService)

	service := NewSASTService(appConfig, mockLinter, mockChecker, mockDockerParser)

	mockParsedDockerfile := &parser.DockerfileParsedResult{
		UsedImages: []string{"alpine:latest"},
	}
	mockDockerParser.On("Parse", mock.Anything).Return(mockParsedDockerfile, nil)

	mockLinterResults := []linter.CheckResult{{}}
	mockLinter.On("Check", mock.Anything).Return(mockLinterResults, nil)

	mockCheckerResults := []vul_checker.VulListItem{{}}
	mockChecker.On("GetCVEByImageNames", mock.Anything).Return(mockCheckerResults)

	// Mock input data
	filename := "Dockerfile"
	content := "FROM alpine:latest"

	// Execute the function
	results := service.Analyze(filename, content)

	// FIXME
	// Verify the mocks
	//mockDockerParser.AssertExpectations(t)
	//mockLinter.AssertExpectations(t)

	// Verify the results
	assert.NotNil(t, results)
}
