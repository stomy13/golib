package main

import (
	"context"
	"log"
	"math/rand"
	"net/http"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	pingv1 "github.com/stomy13/golib/api/internal/gen/connect/ping/v1"
	"github.com/stomy13/golib/api/internal/gen/connect/ping/v1/pingv1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// ランダムな文字列を生成する関数
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func useConnect() {
	mux := http.NewServeMux()

	// CORS設定を追加
	interceptor := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, Connect-Protocol-Version")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

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
		// CORSミドルウェアを適用
		interceptor(h2c.NewHandler(mux, &http2.Server{})),
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
		Text:   generateRandomString(10), // 10文字のランダムな文字列を生成
	})
	res.Header().Set("Some-Other-Header", "hello!")
	return res, nil
}
