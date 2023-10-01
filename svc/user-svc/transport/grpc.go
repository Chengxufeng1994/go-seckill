package transport

import (
	"context"
	"github.com/Chengxufeng1994/go-seckill/pb"
	"github.com/Chengxufeng1994/go-seckill/svc/user-svc/endpoint"
	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	check grpc.Handler
	pb.UnimplementedUserServiceServer
}

func (srv *grpcServer) Check(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	_, resp, err := srv.check.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserResponse), nil
}

func NewGrpcServer(_ context.Context, endpts endpoint.UserEndpoints, serverTracer grpc.ServerOption) pb.UserServiceServer {
	return &grpcServer{
		check: grpc.NewServer(
			endpts.UserEndpoint,
			DecodeGRPCUserRequest,
			EncodeGRPCUserResponse,
			serverTracer,
		),
	}
}

func DecodeGRPCUserRequest(ctx context.Context, request any) (any, error) {
	req := request.(*pb.UserRequest)
	return endpoint.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func EncodeGRPCUserResponse(_ context.Context, response any) (any, error) {
	resp := response.(endpoint.UserResponse)
	if resp.Error != "" {
		return &pb.UserResponse{
			Result: bool(resp.Result),
			UserId: resp.UserId,
			Err:    "error",
		}, nil
	}

	return &pb.UserResponse{
		Result: bool(resp.Result),
		UserId: resp.UserId,
		Err:    "",
	}, nil
}
