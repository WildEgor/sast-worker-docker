package linter

import (
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"os"
)

var _ ILinter = (*TrivyAdapter)(nil)

type TrivyAdapter struct {
	appConfig *configs.AppConfig
}

func NewTrivyAdapter(appConfig *configs.AppConfig) *TrivyAdapter {
	return &TrivyAdapter{
		appConfig,
	}
}

func (l *TrivyAdapter) Check(file *os.File) ([]CheckResult, error) {
	cr := make([]CheckResult, 0)
	// TODO: impl me
	return cr, nil
}
