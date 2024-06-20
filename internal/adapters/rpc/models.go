package rpc

import "context"

type IRPCServer interface {
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}
