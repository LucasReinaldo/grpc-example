package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/LucasReinaldo/grpc-example/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:5051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Error connecting %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)

}

func AddUser(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Lucas",
		Email: "lucas@lucas.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not send gRPC request %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {
	req := &pb.User{
		Id:    "0",
		Name:  "Lucas",
		Email: "lucas@lucas.com",
	}

	resStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not send gRPC request %v", err)
	}

	for {
		stream, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not receive message: %v", err)
		}

		fmt.Println("Status: ", stream.Status, "-", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		{
			Id:    "0",
			Name:  "Lucas Lucas",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "1",
			Name:  "Lucas Lucas1",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "2",
			Name:  "Lucas Lucas2",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "3",
			Name:  "Lucas Lucas3",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "4",
			Name:  "Lucas Lucas4",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "5",
			Name:  "Lucas Lucas5",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "6",
			Name:  "Lucas Lucas6",
			Email: "Lucas2@lucas.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Could not send users, %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error in response, %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {

	reqs := []*pb.User{
		{
			Id:    "0",
			Name:  "Lucas Lucas",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "1",
			Name:  "Lucas Lucas1",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "2",
			Name:  "Lucas Lucas2",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "3",
			Name:  "Lucas Lucas3",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "4",
			Name:  "Lucas Lucas4",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "5",
			Name:  "Lucas Lucas5",
			Email: "Lucas2@lucas.com",
		},
		{
			Id:    "6",
			Name:  "Lucas Lucas6",
			Email: "Lucas2@lucas.com",
		},
	}

	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Could not send users, %v", err)
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			stream.Send(req)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error receiving data, %v", err)
				break
			}
			fmt.Printf("Receiving user %v, status %v \n", res.GetUser().GetName(), res.GetStatus())
		}
		close(wait)
	}()

	<-wait
}
