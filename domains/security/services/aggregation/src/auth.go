package src

import (
	"context"

	pb "github.com/kirvader/BodyController/domains/security/services/aggregation/proto"
	pbAuth "github.com/kirvader/BodyController/domains/security/services/base/auth/proto"
)

func (svc *SecurityService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	_, err := svc.authClient.CreateUser(ctx, &pbAuth.CreateUserRequest{
		UserCredentials: req.UserCredentials,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateUserResponse{}, nil
}

func (svc *SecurityService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	_, err := svc.authClient.DeleteUser(ctx, &pbAuth.DeleteUserRequest{
		Username: req.Username,
	})
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, nil
}

func (svc *SecurityService) LogIn(ctx context.Context, req *pb.LogInRequest) (*pb.LogInResponse, error) {
	resp, err := svc.authClient.LogIn(ctx, &pbAuth.LogInRequest{
		UserCredentials: req.UserCredentials,
	})
	if err != nil {
		return nil, err
	}
	return &pb.LogInResponse{
		Token: resp.Token,
	}, nil
}
