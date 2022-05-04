package main

import (
	"context"
	"testing"

	"google.golang.org/grpc/status"

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

// Test the PostGraph function
func TestGraphServer_PostGraph(t *testing.T) {
	// Post two graphs to the server, one is valid,
	// and the other is not
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
			"",
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
	s := newServer()

	// Check the output of the function
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			id, err := s.PostGraph(ctx, tt.graph)

			if id != nil {
				if id.Id != tt.res.Id {
					t.Error("response: expected", tt.res.Id, "received", id.Id)
				}
			}

			if err != nil {
				if er, ok := status.FromError(err); ok {
					if er.Message() != tt.errMsg {
						t.Error("error message: expected", tt.errMsg, "received", er.Message())
					}
				}
			}
		})
	}
}

// Test the ShortestPath function
func TestGraphServer_ShortestPath(t *testing.T) {

	ctx := context.Background()
	s := newServer()

	// First, post a graph to the server
	g := &pb.Graph{Vertices: []int32{1, 2, 3, 4, 5},
		Edges: map[int32]*pb.Neighbors{
			1: {Neighbors: []int32{2, 3, 4, 5}},
			2: {Neighbors: []int32{1}},
			3: {Neighbors: []int32{1}},
			4: {Neighbors: []int32{1}},
			5: {Neighbors: []int32{1}},
		}}

	id, err0 := s.PostGraph(ctx, g)
	if err0 != nil {
		t.Error("cannot post graph", err0)
	}

	// Query for 2 shortest paths
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
			"",
		},
		{
			"longer path",
			&pb.PathRequest{S: 3, T: 2, Gid: id},
			&pb.Path{Path: []int32{3, 1, 2}},
			"found edge between non-existant nodes",
		},
	}

	// Check the output of the function
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			path, err := s.ShortestPath(ctx, tt.req)

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
		})
	}
}

// Test the DeleteGraph function
func TestGraphServer_DeleteGraph(t *testing.T) {

	// First post a graph, and then execute 2 delete operations:
	// One is valid and the other is not
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
			"",
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
			"non-existant graph",
		},
	}

	ctx := context.Background()
	s := newServer()

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// First check the output of posting the graph
			if i == 0 {
				id, err := s.PostGraph(ctx, tt.graph)

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
				// Check the output of deleting the graphs
			} else {

				rep, err := s.DeleteGraph(ctx, tt.id)

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
