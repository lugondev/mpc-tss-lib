package common

import (
	"context"

	"github.com/ambrosus/ambrosus-bridge/relay/pkg/mpc/tss_wrap"
)

var (
	CenterOperation    = []byte("center")
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