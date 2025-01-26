package main

import (
	"context"
	"log"
	"net"
	pb "github.com/sleektea/teashop/teashop_proto"
	grpc "google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTeaShopServer
}

func (s *server) GetMenu(menuRequest *pb.MenuRequest, srv pb.TeaShop_GetMenuServer) error {
	itemprice := []*pb.ItemPrice{
		&pb.ItemPrice{
			Id: "1",
			Name: "Ginger Tea",
			Price: "20",
		},
		&pb.ItemPrice{
			Id: "2",
			Name: "Masala Tea",
			Price: "25",
		},
		&pb.ItemPrice{
			Id: "3",
			Name: "Samosa",
			Price: "50",
		},
	}
	log.Printf("Itemprice %v", itemprice)
	for i := range itemprice {
		err := srv.Send(&pb.Menu{
			Itemprice: itemprice[0:i+1],
		})
		if err != nil {
			log.Fatalf("Error sending items")
		}
	}
	return nil
}
func (s *server) PlaceOrder(context.Context, *pb.Order) (*pb.Receipt, error) {
	return &pb.Receipt{
		Id: "ABC123",
		Cost: 120,
	}, nil
}
func (s *server) GetOrderStatus(ctx context.Context, receipt *pb.Receipt) (*pb.OrderStatus, error) {
	return &pb.OrderStatus{
		OrderId: receipt.Id,
		Status: "IN PROGRESS",
	}, nil
}
func main() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTeaShopServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis) ; err != nil {
		log.Fatalf("failed to server %s", err)
	}
}
