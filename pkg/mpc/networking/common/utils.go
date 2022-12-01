package common

import (
	"context"

	"github.com/lugondev/mpc-tss-lib/pkg/mpc/tss_wrap"
)

var (
	GatewayOperation   = []byte("gateway")
	KeygenOperation    = []byte("keygen")
	ReshareOperation   = []byte("reshare")
	HeaderTssID        = "X-TSS-ID"
	HeaderTssOperation = "X-TSS-Operation"
	HeaderAccessToken  = "X-Access-Token"
	ResultPrefix       = []byte("result")
)

type OpError struct {
	Type string
	Err  error
}

type OperationFunc func(ctx context.Context, inCh <-chan []byte, outCh chan<- *tss_wrap.Message) ([]byte, error)
