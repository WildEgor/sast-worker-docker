package parser

import (
	"os"
)

// DockerfileParsedResult save parsed data from Dockefile
type DockerfileParsedResult struct {
	UsedImages []string // save image:tag from instruction FROM <image:tag>
}

// IParser for Dockerfile parsers
type IParser interface {
	Parse(file *os.File) (*DockerfileParsedResult, error)
}
