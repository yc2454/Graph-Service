package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/yc2454/Graph-Service/graph_service"
)

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGraphServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	g := new(pb.Graph)
	g.Vertices = []int32{1, 2, 3, 4, 5}
	g.Edges = make(map[int32]*pb.Neighbors)

	var nn [5]*pb.Neighbors
	for i := 0; i < 5; i++ {
		nn[i] = new(pb.Neighbors)
	}
	nn[0].Neighbors = []int32{2, 3, 4, 5}
	nn[1].Neighbors = []int32{1}
	nn[2].Neighbors = []int32{1}
	nn[3].Neighbors = []int32{1}
	nn[4].Neighbors = []int32{1}

	for i := 1; i < 6; i++ {
		g.Edges[int32(i)] = nn[i-1]
	}

	log.Printf("Posting graph")
	id, err := c.PostGraph(ctx, g)

	if err != nil {
		log.Fatalf("could not post graph: %v", err)
	}

	log.Printf("Got ID: %v", id.Id)

	req := new(pb.PathRequest)
	req.Gid = id
	req.S = 2
	req.T = 3

	log.Printf("Asking for the shortest path between %v and %v", req.S, req.T)
	path, err1 := c.ShortestPath(ctx, req)
	log.Printf("Received reply for shortest path")

	if err1 != nil {
		log.Fatalf("could find shortest graph: %v", err1)
	}

	fmt.Printf("The shortest path between %v and %v is: ", req.S, req.T)
	if path != nil {
		for _, n := range path.Path {
			fmt.Printf("%v ", n)
		}
		fmt.Println()
	}

	log.Printf("Deleting graph %v", id.Id)
	reply, err2 := c.DeleteGraph(ctx, id)
	if err2 != nil {
		log.Fatalf("could find delete graph: %v", err1)
	}

	log.Printf(reply.Result)
}
