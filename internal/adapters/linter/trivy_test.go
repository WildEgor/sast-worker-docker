package linter

import (
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestTrivyAdapter_Check(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "temp")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file to represent the Dockerfile
	tempFile, err := os.CreateTemp("", "Dockerfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some content to the Dockerfile
	if _, err := tempFile.WriteString("FROM ubuntu:18.04\n"); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	appConfig := &configs.AppConfig{
		ScriptsPath: "/scripts",
	}
	adapter := NewTrivyAdapter(appConfig)

	results, err := adapter.Check(tempFile)

	assert.NoError(t, err)
	assert.NotNil(t, results)
}
