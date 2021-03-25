package services

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/LucasReinaldo/grpc-example/pb"
	"github.com/google/uuid"
)

// type UserServiceServer interface {
// 	AddUser(context.Context, *User) (*User, error)
// 	mustEmbedUnimplementedUserServiceServer()
// }

// serve para garantir que se adicionarmos um service novo o servidor não vai cair
type UserService struct {
	pb.UnimplementedUserServiceServer
}

func (*UserService) AddUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	// insert in database
	fmt.Println(req.GetName())

	return &pb.User{
		Id:    uuid.New().String(),
		Name:  req.GetName(),
		Email: req.GetEmail(),
	}, nil
}

func (*UserService) AddUserVerbose(req *pb.User, stream pb.UserService_AddUserVerboseServer) error {
	stream.Send(&pb.UserResultStream{
		Status: "Init",
		User:   &pb.User{},
	})

	stream.Send(&pb.UserResultStream{
		Status: "Inserting",
		User:   &pb.User{},
	})

	stream.Send(&pb.UserResultStream{
		Status: "Inserted",
		User: &pb.User{
			Id:    "0",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	stream.Send(&pb.UserResultStream{
		Status: "Completed",
		User: &pb.User{
			Id:    "0",
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	})

	return nil
}

func (*UserService) AddUsers(stream pb.UserService_AddUsersServer) error {
	users := []*pb.User{}

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&pb.Users{
				User: users,
			})
		}

		if err != nil {
			log.Fatalf("Could not send gRPC request %v", err)
		}

		users = append(users, &pb.User{
			Id:    uuid.New().String(),
			Name:  req.GetName(),
			Email: req.GetEmail(),
		})
		fmt.Println(req.GetName())
	}
}

func (*UserService) AddUserStreamBoth(stream pb.UserService_AddUserStreamBothServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			log.Fatalf("Error receiving stream from client: %v", err)
		}

		err = stream.Send(&pb.UserResultStream{
			Status: "Added",
			User:   req,
		})
		if err != nil {
			log.Fatalf("Error sending stream to client: %v", err)
		}
	}
}
