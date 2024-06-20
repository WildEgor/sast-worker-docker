package parser

import (
	"bufio"
	"log/slog"
	"os"
	"strings"
)

var _ IParser = (*Service)(nil)

type Service struct {
}

func NewDockerParserService() *Service {
	return &Service{}
}

// FIXME:
func (s *Service) Parse(file *os.File) (*DockerfileParsedResult, error) {
	parsed := &DockerfileParsedResult{
		UsedImages: make([]string, 0),
	}

	slog.Info("parse dockerfile")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		slog.Info("docker line", slog.Any("value", line))

		lineTokens := strings.Split(line, " ")

		if len(lineTokens) == 0 {
			continue
		}

		switch lineTokens[0] {
		case "FROM":
			if len(lineTokens) == 1 {
				continue
			}
			slog.Info("IMAGE", slog.Any("value", lineTokens[1]))
			parsed.UsedImages = append(parsed.UsedImages, lineTokens[1])
		default:
		}
	}

	if err := scanner.Err(); err != nil {
		slog.Error("error scanning Dockerfile", slog.Any("error", err))
		return nil, err
	}

	slog.Info("parsed dockerfile", slog.Any("value", parsed))

	return parsed, nil
}
