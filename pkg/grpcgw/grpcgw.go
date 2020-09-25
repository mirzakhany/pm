package grpcgw

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"proj/pkg/config"
	"proj/pkg/log"
	"sync"
	"time"
)

type Controller interface {
	InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux)
	InitGrpc(ctx context.Context, server *grpc.Server)
}

type Interceptor struct {
	Unary  grpc.UnaryServerInterceptor
	Stream grpc.StreamServerInterceptor
}

var (
	controllers  []Controller
	interceptors []Interceptor
	lock         sync.RWMutex

	httpPort config.Int
	grpcPort config.Int

	httpAddr string
	grpcAddr string
)

func RegisterController(c Controller) {
	lock.Lock()
	defer lock.Unlock()
	controllers = append(controllers, c)
}

func RegisterInterceptor(i Interceptor) {
	lock.Lock()
	defer lock.Unlock()
	interceptors = append(interceptors, i)
}

func newGrpcServer() *grpc.Server {
	unaryMiddlewares := []grpc.UnaryServerInterceptor{
		grpc_recovery.UnaryServerInterceptor(),
		grpc_ctxtags.UnaryServerInterceptor(),
		grpc_zap.UnaryServerInterceptor(log.Logger()),
	}

	streamMiddlewares := []grpc.StreamServerInterceptor{
		grpc_recovery.StreamServerInterceptor(),
		grpc_ctxtags.StreamServerInterceptor(),
		grpc_zap.StreamServerInterceptor(log.Logger()),
	}

	for i := range interceptors {
		if interceptors[i].Unary != nil {
			unaryMiddlewares = append(unaryMiddlewares, interceptors[i].Unary)
		}
		if interceptors[i].Stream != nil {
			streamMiddlewares = append(streamMiddlewares, interceptors[i].Stream)
		}
	}
	c := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryMiddlewares...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamMiddlewares...)),
	)

	return c
}

// gRPCClient creates a new GRPC client conn
func gRPCClient() (*grpc.ClientConn, error) {
	return grpc.Dial(grpcAddr, grpc.WithInsecure())
}

// Serve start the server and wait
func serveHTTP(ctx context.Context) (func() error, error) {
	var (
		normalMux = http.NewServeMux()
		mux       = runtime.NewServeMux()
	)
	c, err := gRPCClient()
	if err != nil {
		return nil, err
	}

	//normalMux.HandleFunc("/v1/swagger/", swaggerHandler)
	for i := range controllers {
		controllers[i].InitRest(ctx, c, mux)
	}

	normalMux.Handle("/", cors.AllowAll().Handler(mux))
	srv := http.Server{
		Addr:    httpAddr,
		Handler: normalMux,
	}
	log.Info("start http on", log.Any("address", httpAddr))
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	return func() error {
		nCtx, cnl := context.WithTimeout(context.Background(), time.Second)
		defer cnl()

		return srv.Shutdown(nCtx)
	}, nil
}

// Serve start the server and wait
func serveGRPC(ctx context.Context) (func() error, error) {
	srv := newGrpcServer()
	for i := range controllers {
		controllers[i].InitGrpc(ctx, srv)
	}

	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, err
	}
	log.Info("start grpc on", log.Any("address", grpcAddr))
	go func() {
		err := srv.Serve(lis)
		if err != nil {
			log.Error("Connection Closed", log.Err(err))
		}
	}()

	return lis.Close, nil
}

func Serve(ctx context.Context) error {
	lock.RLock()
	defer lock.RUnlock()

	grpcFn, err := serveGRPC(ctx)
	if err != nil {
		return err
	}
	httpFn, err := serveHTTP(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()
	e1 := httpFn()
	e2 := grpcFn()

	if e1 != nil {
		return e1
	}

	return e2
}

func Init() error {
	defaultHttpPort := 8089
	defaultGrpcPort := 9090

	httpPort = config.RegisterInt("server.httpPort", defaultHttpPort)
	grpcPort = config.RegisterInt("server.grpcPort", defaultGrpcPort)

	if err := config.Load(); err != nil {
		log.Panic("load server settings failed")
		return err
	}

	httpAddr = fmt.Sprintf(":%d", httpPort.Int())
	grpcAddr = fmt.Sprintf(":%d", grpcPort.Int())

	return nil
}
