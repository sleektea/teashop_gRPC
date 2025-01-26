package main

import (
	"context"
	"io"
	"log"
	"time"
	pb "github.com/sleektea/teashop/teashop_proto"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:9001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server %v", err)
	}
	defer conn.Close()
	c := pb.NewTeaShopClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	menuStream, err := c.GetMenu(ctx, &pb.MenuRequest{})
	if err != nil {
		log.Fatal("error calling GetMenu")
	}
	done := make(chan bool)
	var itemprice []*pb.ItemPrice

	go func() {
		for {
			resp, err := menuStream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("Cannot recv from Menustream %v", err)
			}
			itemprice = resp.Itemprice
			log.Printf("Item received %v", itemprice)
		}
	}()
	
	select {
	case <- done:
	case <- time.After(5* time.Second) :
		log.Fatalf("Timeout") 
	}

	var itemList []*pb.Item
	for _, i := range itemprice {
		itemList = append(itemList, &pb.Item{Id: i.Id, Name: i.Name,})
	}
	receipt, _ := c.PlaceOrder(ctx, &pb.Order{Items: itemList})
	log.Printf("Receipt : %v ", receipt)
	status, _ := c.GetOrderStatus(ctx, receipt)
	log.Printf("Status : %v ", status)

}