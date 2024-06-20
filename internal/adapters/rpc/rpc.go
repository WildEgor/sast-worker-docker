package rpc

import (
	"context"
	"fmt"
	"github.com/WildEgor/sast-worker-docker/internal/configs"
	"github.com/WildEgor/sast-worker-docker/internal/services/sast"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log/slog"
	"net"
)

var _ IRPCServer = (*ServerAdapter)(nil)

// FIXME: i guess its more clear to make /router dir with Router that call gRPC server, but there are some limitations - generated pb files must be in same package

// ServerAdapter represent adapter for gRPC server
type ServerAdapter struct {
	appConfig   *configs.AppConfig
	sastService *sast.Service
	srv         *grpc.Server
}

func NewRPCServerAdapter(
	cfg *configs.AppConfig,
	ss *sast.Service,
) *ServerAdapter {
	return &ServerAdapter{
		cfg,
		ss,
		nil,
	}
}

// Run gRPC server
func (s *ServerAdapter) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.appConfig.GRPCPort))
	if err != nil {
		return err
	}

	go func() {
		s.srv = grpc.NewServer()
		RegisterSASTWorkerServer(s.srv, s)
		reflection.Register(s.srv)

		select {
		case <-ctx.Done():
			return
		default:
			slog.Info("grpc server listen")

			if err = s.srv.Serve(lis); err != nil {
				slog.Error("error serve grpc", slog.Any("err", err))
				panic(err)
			}
		}
	}()

	return nil
}

// Stop gRPC server
func (s *ServerAdapter) Stop(ctx context.Context) error {
	s.srv.Stop()
	return nil
}

// ROUTING

// HealthCheck just act like ping
func (s *ServerAdapter) HealthCheck(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	return nil, nil
}

// AnalyzeFile return file analyze
func (s *ServerAdapter) AnalyzeFile(ctx context.Context, request *AnalyzeDockerfileRequest) (*AnalyzeFileResponse, error) {
	r := s.sastService.Analyze(request.File.Filename, request.File.Content)

	slog.Info("result", slog.Any("value", r))

	res := &AnalyzeFileResponse{
		Result: make([]*AnalyzeFileResult, 0),
	}

	if len(r) != 0 {
		for _, result := range r {
			res.Errors += 1
			res.Result = append(res.Result, &AnalyzeFileResult{
				Line: result.Line,
				Coll: result.Coll,
				Msg:  result.Msg,
				Code: result.Code,
			})
		}
	}

	return res, nil
}

func (s *ServerAdapter) mustEmbedUnimplementedSASTWorkerServer() {
}
