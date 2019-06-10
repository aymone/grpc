package handler

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/aymone/grpc/proto/service"
)

type handler struct {
	users map[string]pb.User
}

// New Handler
func New() handler {
	return handler{
		users: make(map[string]pb.User),
	}
}

func (s handler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*empty.Empty, error) {
	log.Println("Creating user...")
	user := req.GetUser()

	if user.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	if user.Role == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "role cannot be empty")
	}

	s.users[user.Username] = *user

	log.Println("User created!")
	return &empty.Empty{}, nil
}

func (s handler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	log.Println("Getting user!")

	if req.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	u, exists := s.users[req.Username]
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}

	log.Println("User found!")
	return &u, nil
}

func (s handler) GreetUser(ctx context.Context, req *pb.GreetUserRequest) (*pb.GreetUserResponse, error) {
	log.Println("Greeting user...")
	if req.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}
	if req.Greeting == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "greeting cannot be empty")
	}

	user, err := s.GetUser(ctx, &pb.GetUserRequest{Username: req.Username})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find matching user")
	}

	return &pb.GreetUserResponse{
		Greeting: fmt.Sprintf("%s, %s! You are a great %s!", strings.Title(req.Greeting), user.Username, user.Role),
	}, nil
}
