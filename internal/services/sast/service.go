package sast

import (
	"fmt"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/linter"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/vul_checker"
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"github.com/WildEgor/sast-worker-docker/internal/services/parser"
	"github.com/WildEgor/sast-worker-docker/internal/utils"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"path"
)

// Service represent SAST logic for Dockerfiles
type Service struct {
	appConfig    *configs.AppConfig
	linter       linter.ILinter
	vulChecker   vul_checker.IVulChecker
	dockerParser parser.IParser
}

// NewSASTService create instance
func NewSASTService(appConfig *configs.AppConfig, linter linter.ILinter, vulChecker vul_checker.IVulChecker, dockerParser parser.IParser) *Service {
	return &Service{
		appConfig,
		linter,
		vulChecker,
		dockerParser,
	}
}

// Analyze run Dockerfile analysis process and return result where content is string representation of Dockerfile
func (s *Service) Analyze(filename string, content string) []AnalysisItem {
	tmpFile, err := s.saveTemp(filename, content)
	if err != nil {
		slog.Error("save file fail", slog.Any("err", err))
		return nil
	}
	defer s.removeTemp(tmpFile)

	checkDockerfileConfig := func(file *os.File) <-chan AnalysisItem {
		out := make(chan AnalysisItem)

		go func() {
			defer close(out)

			check, err := s.linter.Check(file)
			if err != nil {
				slog.Error("linter err", slog.Any("err", err))
				return
			}

			for _, result := range MapCheckResultToAnalysisResult(check) {
				out <- result
			}
		}()

		return out
	}

	checkDockerfileVuls := func(file *os.File) <-chan AnalysisItem {
		out := make(chan AnalysisItem)

		go func() {
			defer close(out)

			slog.Info("try parse dockerfile")

			parsed, err := s.dockerParser.Parse(file)
			if err != nil {
				slog.Error("docker parse err", slog.Any("err", err))
				return
			}

			vuls := s.vulChecker.GetCVEByImageNames(parsed.UsedImages)
			for _, result := range MapVulItemsToAnalysisResult(vuls) {
				out <- result
			}
		}()

		return out
	}

	outR := make([]AnalysisItem, 0)
	result := utils.ChannelMerger[AnalysisItem](checkDockerfileVuls(tmpFile), checkDockerfileConfig(tmpFile))
	for r := range result {
		outR = append(outR, r)
	}

	return outR
}

func (s *Service) saveTemp(filename, content string) (*os.File, error) {
	fpath := path.Join(s.appConfig.TempPath, fmt.Sprintf("%s_%s", filename, uuid.NewString()))

	tmpFile, err := os.Create(fpath)
	if err != nil {
		slog.Error("create temp file fail", slog.Any("err", err))
		return nil, err
	}

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		slog.Error("write temp file fail", slog.Any("err", err))
		return nil, err
	}

	if _, err = tmpFile.Seek(0, 0); err != nil {
		slog.Error("seek temp file fail", slog.Any("err", err))
		return nil, err
	}

	return tmpFile, nil
}

func (s *Service) removeTemp(file *os.File) error {
	file.Close()
	return os.Remove(file.Name())
}
