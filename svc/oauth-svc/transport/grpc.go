package transport

import (
	"context"
	"github.com/Chengxufeng1994/go-seckill/pb"
	"github.com/Chengxufeng1994/go-seckill/svc/oauth-svc/endpoint"
	"github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	checkToken grpc.Handler
	pb.UnimplementedOAuthServiceServer
}

func (s *grpcServer) CheckToken(ctx context.Context, r *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, resp, err := s.checkToken.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CheckTokenResponse), nil
}

func NewGrpcServer(_ context.Context, endpts endpoint.OAuth2Endpoints, serverTracer grpc.ServerOption) pb.OAuthServiceServer {
	return &grpcServer{
		checkToken: grpc.NewServer(
			endpts.GrpcCheckTokenEndpoint,
			DecodeGRPCCheckTokenRequest,
			EncodeGRPCCheckTokenResponse,
			serverTracer,
		),
	}
}

func DecodeGRPCCheckTokenRequest(ctx context.Context, request any) (any, error) {
	req := request.(*pb.CheckTokenRequest)
	return &endpoint.CheckTokenRequest{
		Token: req.Token,
	}, nil
}

func EncodeGRPCCheckTokenResponse(_ context.Context, response any) (any, error) {
	resp := response.(*endpoint.CheckTokenResponse)
	if resp.Error != "" {
		return &pb.CheckTokenResponse{
			IsValidToken: false,
			Err:          "error",
		}, nil
	}

	return &pb.CheckTokenResponse{
		UserDetails: &pb.UserDetails{
			UserId:      string(resp.OAuthDetails.User.UserId),
			Username:    resp.OAuthDetails.User.Username,
			Authorities: resp.OAuthDetails.User.Authorities,
		},
		ClientDetails: &pb.ClientDetails{
			ClientId:                    resp.OAuthDetails.Client.ClientId,
			AccessTokenValiditySeconds:  int32(resp.OAuthDetails.Client.AccessTokenValiditySeconds),
			RefreshTokenValiditySeconds: int32(resp.OAuthDetails.Client.RefreshTokenValiditySeconds),
			AuthorizedGrantTypes:        resp.OAuthDetails.Client.AuthorizedGrantTypes,
		},
		IsValidToken: true,
		Err:          "",
	}, nil
}
