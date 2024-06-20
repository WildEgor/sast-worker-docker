package linter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

var _ ILinter = (*HadolintAdapter)(nil)

// HadolintAdapter wrapper for hadolint CLI
type HadolintAdapter struct {
	appConfig *configs.AppConfig
}

func NewHadolintAdapter(appConfig *configs.AppConfig) *HadolintAdapter {
	return &HadolintAdapter{
		appConfig,
	}
}

// Check call hadolint CLI to check Dockerfile
func (l *HadolintAdapter) Check(file *os.File) ([]CheckResult, error) {
	cr := make([]CheckResult, 0)
	cargs := make([]string, 0)

	cargs = append(cargs, fmt.Sprintf("%s/hadolint/run.sh", l.appConfig.ScriptsPath))
	cargs = append(cargs, file.Name())

	slog.Info("hadolint command args", slog.Any("value", strings.Join(cargs, " ")))

	r, err := exec.Command("/bin/bash", cargs[0], cargs[1]).CombinedOutput()
	if err != nil {
		var exitError *exec.ExitError
		ok := errors.As(err, &exitError)
		if !ok {
			slog.Error("check hadolint err", slog.Any("err", err))
			return cr, err
		}
	}

	lr := make([]HadolintResult, 0)
	if err := json.Unmarshal(r, &lr); err != nil {
		slog.Error("unmarshal hadolint result err", slog.Any("err", err))
		return cr, err
	}

	result := make([]CheckResult, len(lr))
	for i, hr := range lr {
		result[i] = CheckResult{
			Pos: CheckPos{
				Coll: hr.Coll,
				Line: hr.Line,
			},
			Err: CheckError{
				Code:  hr.Code,
				Level: hr.Level,
				Msg:   hr.Message,
			},
		}
	}

	return result, nil
}
