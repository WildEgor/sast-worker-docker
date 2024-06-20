package parser

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempFile(t *testing.T, content string) *os.File {
	file, err := os.CreateTemp("", "dockerfile-*.Dockerfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to seek to beginning of temp file: %v", err)
	}

	return file
}

func TestService_Parse(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    *DockerfileParsedResult
		expectError bool
	}{
		{
			name:        "valid Dockerfile",
			content:     "FROM golang:1.16\nFROM alpine:3.13\n",
			expected:    &DockerfileParsedResult{UsedImages: []string{"golang:1.16", "alpine:3.13"}},
			expectError: false,
		},
		{
			name:        "empty Dockerfile",
			content:     "",
			expected:    &DockerfileParsedResult{UsedImages: []string{}},
			expectError: false,
		},
		{
			name:        "Dockerfile with missing image name",
			content:     "FROM\n",
			expected:    &DockerfileParsedResult{UsedImages: []string{}},
			expectError: false,
		},
		{
			name:        "Dockerfile with invalid instruction",
			content:     "RUN echo 'hello world'\n",
			expected:    &DockerfileParsedResult{UsedImages: []string{}},
			expectError: false,
		},
	}

	service := NewDockerParserService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := createTempFile(t, tt.content)
			defer os.Remove(file.Name())
			defer file.Close()

			result, err := service.Parse(file)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
