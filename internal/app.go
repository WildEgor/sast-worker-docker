package app

import (
	"context"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/linter"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/rpc"
	"github.com/WildEgor/sast-worker-docker/internal/adapters/vul_checker"
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"github.com/WildEgor/sast-worker-docker/internal/services/parser"
	"github.com/WildEgor/sast-worker-docker/internal/services/sast"
	"log/slog"
	"os"
)

// App represents the main server configuration.
type App struct {
	RPC *rpc.ServerAdapter
}

// NewApp initialize all deps
func NewApp() (*App, error) {
	// Init configs
	configs.Init()
	configs.NewLoggerConfig()
	appConfig := configs.NewAppConfig()
	trivyConfig := configs.NewTrivyConfig()

	// Init adapters
	linterAdapter := linter.NewHadolintAdapter(appConfig)
	trivyChecker := vul_checker.NewTrivyCheckerAdapter(appConfig, trivyConfig)
	dockerParser := parser.NewDockerParserService()

	// Init services
	sastService := sast.NewSASTService(appConfig, linterAdapter, trivyChecker, dockerParser)

	return &App{
		RPC: rpc.NewRPCServerAdapter(appConfig, sastService),
	}, nil
}

// Run start service with deps
func (srv *App) Run(ctx context.Context) {
	slog.Info("server is listening")

	if err := srv.RPC.Run(ctx); err != nil {
		slog.Error("server start fail", slog.Any("err", err))
		os.Exit(1)
	}
}

// Shutdown graceful shutdown
func (srv *App) Shutdown(ctx context.Context) {
	slog.Info("shutdown service")

	if err := srv.RPC.Stop(ctx); err != nil {
		slog.Error("shutdown rpc fail", slog.Any("err", err))
		return
	}
}
