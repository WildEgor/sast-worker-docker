package vul_checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"log/slog"
	"os/exec"
	"strings"
)

var _ IVulChecker = (*TrivyCheckerAdapter)(nil)

type TrivyCheckerAdapter struct {
	appConfig   *configs.AppConfig
	trivyConfig *configs.TrivyConfig
}

func NewTrivyCheckerAdapter(appConfig *configs.AppConfig, trivyConfig *configs.TrivyConfig) *TrivyCheckerAdapter {
	return &TrivyCheckerAdapter{
		appConfig,
		trivyConfig,
	}
}

func (c *TrivyCheckerAdapter) GetCVEByImageNames(names []string) []VulListItem {
	vi := make([]VulListItem, 0)

	slog.Info("try get cve", slog.Any("names", names))

	for _, name := range names {
		slog.Info("check cve in image", slog.Any("value", name))

		cargs := make([]string, 0)
		cargs = append(cargs, fmt.Sprintf("%s/trivy/client.sh", c.appConfig.ScriptsPath)) // 0
		cargs = append(cargs, "image")
		cargs = append(cargs, c.trivyConfig.API)
		cargs = append(cargs, name)

		slog.Info("trivy checker command args", strings.Join(cargs, " "))

		r, err := exec.Command("/bin/bash", cargs[0], cargs[1], cargs[2], cargs[3]).CombinedOutput()
		if err != nil {
			var exitError *exec.ExitError
			ok := errors.As(err, &exitError)
			if !ok {
				slog.Error("check hadolint err", slog.Any("err", err))
				return vi
			}
		}

		tr := &TrivyImageCheckResult{}
		if err := json.Unmarshal(r, tr); err != nil {
			slog.Error("unmarshal trivy checker result err", slog.Any("err", err))
			continue // TODO: specify domain lvl error
		}

		slog.Info("result", slog.Any("value", tr))

		for _, result := range tr.Results {
			for _, vulnerability := range result.Vulnerabilities {
				vi = append(vi, VulListItem{
					Code: vulnerability.VulnerabilityID,
					Msg:  vulnerability.Desc,
				})
			}
		}
	}

	return vi
}
