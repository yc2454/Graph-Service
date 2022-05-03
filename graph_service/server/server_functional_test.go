package graphservice

import (
	"context"
	"fmt"
	"log"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/yc2454/Graph-Service/graph_service"
)

func dialer() func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()

	pb.RegisterGraphServiceServer(server, newServer())

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

func TestGraphServer_SingleClient(t *testing.T) {
	tests := []struct {
		name   string
		graph  *pb.Graph
		id     *pb.GraphID
		msg    *pb.DeleteReply
		req    *pb.PathRequest
		res    *pb.Path
		errMsg string
	}{
		{
			"post a valid graph",
			&pb.Graph{Vertices: []int32{1, 2, 3, 4, 5},
				Edges: map[int32]*pb.Neighbors{
					1: {Neighbors: []int32{2, 3, 4, 5}},
					2: {Neighbors: []int32{1, 3}},
					3: {Neighbors: []int32{1, 2}},
					4: {Neighbors: []int32{1}},
					5: {Neighbors: []int32{1}},
				}},
			&pb.GraphID{Id: 1},
			nil,
			nil,
			nil,
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"find a shortest path",
			nil,
			&pb.GraphID{Id: 1},
			nil,
			&pb.PathRequest{S: 1, T: 2, Gid: &pb.GraphID{Id: 1}},
			&pb.Path{Path: []int32{1, 2}},
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"find another shortest path",
			nil,
			&pb.GraphID{Id: 1},
			nil,
			&pb.PathRequest{S: 2, T: 3, Gid: &pb.GraphID{Id: 1}},
			&pb.Path{Path: []int32{2, 3}},
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"delete prev graph",
			nil,
			&pb.GraphID{Id: 1},
			&pb.DeleteReply{Result: "Successfully deleted the graph"},
			nil,
			nil,
			"found edge between non-existant nodes",
		},
		{
			"delete non-existant graph",
			nil,
			&pb.GraphID{Id: 1},
			&pb.DeleteReply{Result: "No such id stored"},
			nil,
			nil,
			"non-existant graph",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGraphServiceClient(conn)

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if i == 0 {
				id, err := client.PostGraph(ctx, tt.graph)

				if id != nil {
					if id.Id != tt.id.Id {
						t.Error("response: expected", tt.id.Id, "received", id.Id)
					}
				}

				if err != nil {
					if er, ok := status.FromError(err); ok {
						if er.Message() != tt.errMsg {
							t.Error("error message: expected", tt.errMsg, "received", er.Message())
						}
					}
				}
			} else if i == 1 || i == 2 {
				path, err := client.ShortestPath(ctx, tt.req)

				if path != nil {
					if !Equal(path.Path, tt.res.Path) {
						t.Error("response: expected", tt.res.Path, "received", path.Path)
					}
				}

				if err != nil {
					if er, ok := status.FromError(err); ok {
						if er.Message() != tt.errMsg {
							t.Error("error message: expected", tt.errMsg, "received", er.Message())
						}
					}
				}
			} else {

				rep, err := client.DeleteGraph(ctx, tt.id)

				if rep != nil {
					if rep.Result != tt.msg.Result {
						t.Error("response: expected", tt.msg.Result, "received", rep.Result)
					}
				}

				if err != nil {
					if er, ok := status.FromError(err); ok {
						if er.Message() != tt.errMsg {
							t.Error("error message: expected", tt.errMsg, "received", er.Message())
						}
					}
				}
			}
		})
	}
}
