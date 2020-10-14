package issues

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/mirzakhany/pm/pkg/grpcgw"
	"github.com/mirzakhany/pm/protobuf/issues"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type API interface {
	grpcgw.Controller
	issues.IssueServiceServer
}

type api struct {
	service Service
}

func (a api) InitRest(ctx context.Context, conn *grpc.ClientConn, mux *runtime.ServeMux) {
	cl := issues.NewIssueServiceClient(conn)
	_ = issues.RegisterIssueServiceHandlerClient(ctx, mux, cl)
}

func (a api) InitGrpc(ctx context.Context, server *grpc.Server) {
	issues.RegisterIssueServiceServer(server, a)
}

func (a api) ListIssues(ctx context.Context, request *issues.ListIssuesRequest) (*issues.ListIssuesResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.Query(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetIssue(ctx context.Context, request *issues.GetIssueRequest) (*issues.Issue, error) {
	res, err := a.service.Get(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateIssue(ctx context.Context, request *issues.CreateIssueRequest) (*issues.Issue, error) {
	res, err := a.service.Create(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateIssue(ctx context.Context, request *issues.UpdateIssueRequest) (*issues.Issue, error) {
	res, err := a.service.Update(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteIssue(ctx context.Context, request *issues.DeleteIssueRequest) (*empty.Empty, error) {
	_, err := a.service.Delete(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return nil, err
}

func (a api) ListIssueStatus(ctx context.Context, request *issues.ListIssueStatusRequest) (*issues.ListIssueStatusResponse, error) {
	offset, limit := grpcgw.GetOffsetAndLimit(request.Offset, request.Limit)
	res, err := a.service.QueryStatus(ctx, offset, limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) GetIssueStatus(ctx context.Context, request *issues.GetIssueStatusRequest) (*issues.IssueStatus, error) {
	res, err := a.service.GetStatus(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) CreateIssueStatus(ctx context.Context, request *issues.CreateIssueStatusRequest) (*issues.IssueStatus, error) {
	res, err := a.service.CreateStatus(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) UpdateIssueStatus(ctx context.Context, request *issues.UpdateIssueStatusRequest) (*issues.IssueStatus, error) {
	res, err := a.service.UpdateStatus(ctx, request)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return res, err
}

func (a api) DeleteIssueStatus(ctx context.Context, request *issues.DeleteIssueStatusRequest) (*empty.Empty, error) {
	_, err := a.service.DeleteStatus(ctx, request.Uuid)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return nil, err
}

func (a api) SetIssueStatus(ctx context.Context, request *issues.SetIssueStatusRequest) (*issues.Issue, error) {
	panic("implement me")
}

func New(srv Service) API {
	s := api{service: srv}
	grpcgw.RegisterController(s)
	return s
}
