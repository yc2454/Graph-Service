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

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []int32) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

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

func TestGraphServer_PostGraph(t *testing.T) {
	tests := []struct {
		name   string
		graph  *pb.Graph
		res    *pb.GraphID
		errMsg string
	}{
		{
			"valid graph",
			&pb.Graph{Vertices: []int32{1, 2, 3, 4, 5},
				Edges: map[int32]*pb.Neighbors{
					1: {Neighbors: []int32{2, 3, 4, 5}},
					2: {Neighbors: []int32{1}},
					3: {Neighbors: []int32{1}},
					4: {Neighbors: []int32{1}},
					5: {Neighbors: []int32{1}},
				}},
			&pb.GraphID{Id: 1},
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"invalid graph with edge between non-existant nodes",
			&pb.Graph{Vertices: []int32{1, 2, 3, 4, 5},
				Edges: map[int32]*pb.Neighbors{
					1: {Neighbors: []int32{2, 3, 4, 5, 6}},
					2: {Neighbors: []int32{1}},
					3: {Neighbors: []int32{1}},
					4: {Neighbors: []int32{1}},
					5: {Neighbors: []int32{1}},
				}},
			nil,
			"found edge between non-existant nodes",
		},
	}

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGraphServiceClient(conn)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			id, err := client.PostGraph(ctx, tt.graph)

			if id != nil {
				if id.Id != tt.res.Id {
					t.Error("response: expected", tt.res.Id, "received", id.Id)
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					// if er.Code() != tt.errCode {
					// 	t.Error("error code: expected", codes.Unknown, "received", er.Code())
					// }
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}

func TestGraphServer_ShortestPath(t *testing.T) {

	ctx := context.Background()

	conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewGraphServiceClient(conn)
	g := &pb.Graph{Vertices: []int32{1, 2, 3, 4, 5},
		Edges: map[int32]*pb.Neighbors{
			1: {Neighbors: []int32{2, 3, 4, 5}},
			2: {Neighbors: []int32{1}},
			3: {Neighbors: []int32{1}},
			4: {Neighbors: []int32{1}},
			5: {Neighbors: []int32{1}},
		}}

	id, err0 := client.PostGraph(ctx, g)
	if err0 != nil {
		t.Error("cannot post graph", err0)
	}

	tests := []struct {
		name   string
		req    *pb.PathRequest
		res    *pb.Path
		errMsg string
	}{
		{
			"base case",
			&pb.PathRequest{S: 1, T: 2, Gid: id},
			&pb.Path{Path: []int32{1, 2}},
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"longer path",
			&pb.PathRequest{S: 3, T: 2, Gid: id},
			&pb.Path{Path: []int32{3, 1, 2}},
			"found edge between non-existant nodes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			path, err := client.ShortestPath(ctx, tt.req)

			if id != nil {
				if !Equal(path.Path, tt.res.Path) {
					t.Error("response: expected", tt.res.Path, "received", path.Path)
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					// if er.Code() != tt.errCode {
					// 	t.Error("error code: expected", codes.Unknown, "received", er.Code())
					// }
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}

func TestGraphServer_DeleteGraph(t *testing.T) {
	tests := []struct {
		name   string
		graph  *pb.Graph
		id     *pb.GraphID
		msg    *pb.DeleteReply
		errMsg string
	}{
		{
			"post valid graph",
			&pb.Graph{Vertices: []int32{1, 2, 3, 4, 5},
				Edges: map[int32]*pb.Neighbors{
					1: {Neighbors: []int32{2, 3, 4, 5}},
					2: {Neighbors: []int32{1}},
					3: {Neighbors: []int32{1}},
					4: {Neighbors: []int32{1}},
					5: {Neighbors: []int32{1}},
				}},
			&pb.GraphID{Id: 1},
			nil,
			fmt.Sprintf("cannot deposit %v", -1.11),
		},
		{
			"delete prev graph",
			nil,
			&pb.GraphID{Id: 1},
			&pb.DeleteReply{Result: "Successfully deleted the graph"},
			"found edge between non-existant nodes",
		},
		{
			"delete non-existant graph",
			nil,
			&pb.GraphID{Id: 1},
			&pb.DeleteReply{Result: "No such id stored"},
			"found edge between non-existant nodes",
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
						// if er.Code() != tt.errCode {
						// 	t.Error("error code: expected", codes.Unknown, "received", er.Code())
						// }
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
						// if er.Code() != tt.errCode {
						// 	t.Error("error code: expected", codes.Unknown, "received", er.Code())
						// }
						if er.Message() != tt.errMsg {
							t.Error("error message: expected", tt.errMsg, "received", er.Message())
						}
					}
				}
			}
		})
	}
}
