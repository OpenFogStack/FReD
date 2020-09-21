package main

import (
	"context"
	"fmt"
	grpcclient "gitlab.tu-berlin.de/mcc-fred/fred/pkg/externalconnection"
	"google.golang.org/grpc"
	"os"
)

func main() {
	port := 2000
	host := "localhost"

	conn, _ := grpc.Dial(fmt.Sprintf("%s:%d", host, port), grpc.WithInsecure())
	client := grpcclient.NewClientClient(conn)

	client.CreateKeygroup(context.Background(), &grpcclient.CreateKeygroupRequest{Keygroup: "test"})
	client.Update(context.Background(), &grpcclient.UpdateRequest{
		Keygroup: "test",
		Id:       "test-id",
		Data:     "test-data",
	})
	resp, err := client.Read(context.Background(), &grpcclient.ReadRequest{
		Keygroup: "test",
		Id:       "test-id",
	})
	if resp.Data == "test" {
		print("Read successful!")
		os.Exit(0)
	} else {
		print(err)
		os.Exit(1)
	}
}
