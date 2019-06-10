package handler_test

import (
	"context"
	"testing"

	"github.com/aymone/grpc/handler"
	pb "github.com/aymone/grpc/proto/service"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	h := handler.New()

	t.Run("success", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.CreateUserRequest{
			User: &pb.User{
				Username: "John",
				Role:     "soldier",
			},
		}
		e, err := h.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, e)
	})

	t.Run("username error", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.CreateUserRequest{
			User: &pb.User{
				Role: "soldier",
			},
		}
		e, err := h.CreateUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, e)

		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = username cannot be empty")
	})

	t.Run("role error", func(t *testing.T) {
		ctx := context.Background()
		req := &pb.CreateUserRequest{
			User: &pb.User{
				Username: "John",
			},
		}
		e, err := h.CreateUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, e)

		assert.EqualError(t, err, "rpc error: code = InvalidArgument desc = role cannot be empty")
	})
}
