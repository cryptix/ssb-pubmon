package ssb

import (
	"context"
	"fmt"
	"strings"

	"cryptoscope.co/go/muxrpc"
	kitlog "github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type TryHandler struct {
	remoteID string
	log      kitlog.Logger
	Wait     chan struct{}
}

func NewTryHandler(key string, log kitlog.Logger) *TryHandler {
	w := make(chan struct{})
	return &TryHandler{key, log, w}
}

func (h TryHandler) HandleCall(ctx context.Context, req *muxrpc.Request) {
	// TODO: push manifest check into muxrpc
	if req.Type == "" {
		req.Type = "async"
	}

	m := strings.Join(req.Method, ".")
	switch m {
	case "gossip.ping":
		fallthrough

	case "blobs.createWants":
		fallthrough

	case "createHistoryStream":
		h.Wait <- struct{}{}
		req.Stream.Close()

	default:
		h.log.Log("TryHandler", "unhandled call", "method", m, "args", fmt.Sprintf("%+v", req.Args))
		req.Stream.CloseWithError(errors.Errorf("unhandled call"))
	}
}

func (h TryHandler) HandleConnect(ctx context.Context, e muxrpc.Endpoint) {
	/* calling back
	ret, err := e.Async(ctx, "str", []string{"whoami"})
	if err != nil {
		h.log.Log("handleConnect", "whoami", "err", err)
		return
	}
	*/
	h.log.Log("TryHandler", "connect")
}
