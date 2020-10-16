package templates

const ApiTmpl = `
package {{ .Pkg.NamePlural | lower }}

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/protobuf/{{ .Pkg.NamePlural | lower }}"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}ServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := {{ .Pkg.NamePlural | lower }}.New{{ .Pkg.Name }}ServiceClient(conn)
	_ = {{ .Pkg.NamePlural | lower }}.Register{{ .Pkg.Name }}ServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	{{ .Pkg.NamePlural | lower }}.Register{{ .Pkg.Name }}ServiceServer(server, a)
}

func (a api) List{{ .Pkg.NamePlural }}(ctx context.Context, request *{{ .Pkg.NamePlural | lower }}.List{{ .Pkg.NamePlural }}Request) (*{{ .Pkg.NamePlural | lower }}.List{{ .Pkg.NamePlural }}Response, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) Get{{ .Pkg.Name }}(ctx context.Context, request *{{ .Pkg.NamePlural | lower }}.Get{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) Create{{ .Pkg.Name }}(ctx context.Context, request *{{ .Pkg.NamePlural | lower }}.Create{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) Update{{ .Pkg.Name }}(ctx context.Context, request *{{ .Pkg.NamePlural | lower }}.Update{{ .Pkg.Name }}Request) (*{{ .Pkg.NamePlural | lower }}.{{ .Pkg.Name }}, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) Delete{{ .Pkg.Name }}(ctx context.Context, request *{{ .Pkg.NamePlural | lower }}.Delete{{ .Pkg.Name }}Request) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return nil, err
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}

`
