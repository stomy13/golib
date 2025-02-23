package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	pingv1 "github.com/stomy13/golib/api/internal/gen/connect/ping/v1"
	"github.com/stomy13/golib/api/internal/gen/connect/ping/v1/pingv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func useConnect() {

	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		pingv1connect.PingServiceName,
		// protoc-gen-connect-go generates package-level constants
		// for these fully-qualified protobuf service names, so you'd more likely
		// reference userv1.UserServiceName and groupv1.GroupServiceName.
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	// Many tools still expect the older version of the server reflection API, so
	// most servers should mount both handlers.
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	// If you don't need to support HTTP/2 without TLS (h2c), you can drop
	// x/net/http2 and use http.ListenAndServeTLS instead.
	// The generated constructors return a path and a plain net/http
	// handler.
	mux.Handle(pingv1connect.NewPingServiceHandler(&PingServer{}))
	err := http.ListenAndServe(
		"localhost:8080",
		// For gRPC clients, it's convenient to support HTTP/2 without TLS. You can
		// avoid x/net/http2 by using http.ListenAndServeTLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
	log.Fatalf("listen failed: %v", err)
}

type PingServer struct {
	pingv1connect.UnimplementedPingServiceHandler // returns errors from all methods
}

func (ps *PingServer) Ping(
	ctx context.Context,
	req *connect.Request[pingv1.PingRequest],
) (*connect.Response[pingv1.PingResponse], error) {
	// connect.Request and connect.Response give you direct access to headers and
	// trailers. No context-based nonsense!
	log.Println(req.Header().Get("Some-Header"))
	res := connect.NewResponse(&pingv1.PingResponse{
		// req.Msg is a strongly-typed *pingv1.PingRequest, so we can access its
		// fields without type assertions.
		Number: req.Msg.Number,
	})
	res.Header().Set("Some-Other-Header", "hello!")
	return res, nil
}
