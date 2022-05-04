package main

import (
	"context"
	"testing"

	pb "github.com/yc2454/Graph-Service/graph_service"
)

// Performance test for posting three graphs in a row
func BenchmarkGraphServer_PostGraphPerf(b *testing.B) {

	ctx := context.Background()

	s := newServer()

	gs := []*pb.Graph{
		{Vertices: []int32{1, 2, 3, 4, 5},
			Edges: map[int32]*pb.Neighbors{
				1: {Neighbors: []int32{2, 3, 4, 5}},
				2: {Neighbors: []int32{1}},
				3: {Neighbors: []int32{1}},
				4: {Neighbors: []int32{1}},
				5: {Neighbors: []int32{1}},
			}},
		{Vertices: []int32{1, 2, 3, 4, 5},
			Edges: map[int32]*pb.Neighbors{
				1: {Neighbors: []int32{2}},
				2: {Neighbors: []int32{1, 3}},
				3: {Neighbors: []int32{2, 4}},
				4: {Neighbors: []int32{3, 5}},
				5: {Neighbors: []int32{4}},
			}},
		{Vertices: []int32{1, 2, 3, 4, 5, 6, 7},
			Edges: map[int32]*pb.Neighbors{
				1: {Neighbors: []int32{2, 3, 4, 5}},
				2: {Neighbors: []int32{1}},
				3: {Neighbors: []int32{1}},
				4: {Neighbors: []int32{1}},
				5: {Neighbors: []int32{1}},
			}},
	}

	for i, g := range gs {
		id, err := s.PostGraph(ctx, g)

		if err == nil {
			if id.Id != int32(i+1) {
				b.Error("response: expected", i+1, "received", id.Id)
			}
		}
	}

}

// Performance test for posting a graph and finding
// shortest paths
func BenchmarkGraphServer_ShortestPathPerf(b *testing.B) {

	ctx := context.Background()
	s := newServer()

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
		b.Error("cannot post graph", err0)
	}

	reqs := []*pb.PathRequest{
		{S: 1, T: 2, Gid: id},
		{S: 2, T: 3, Gid: id},
		{S: 1, T: 4, Gid: id},
		{S: 5, T: 4, Gid: id},
	}

	paths := []pb.Path{
		{Path: []int32{1, 2}},
		{Path: []int32{2, 1, 3}},
		{Path: []int32{1, 4}},
		{Path: []int32{5, 1, 4}},
	}

	for i, req := range reqs {

		path, err := s.ShortestPath(ctx, req)

		if path != nil {
			if !Equal(path.Path, paths[i].Path) {
				b.Error("response: expected", paths[i].Path, "received", path.Path)
			}
		}

		if err != nil {
			b.Error(err)
		}
	}

}

// Performance test for posting and deleting a graph
func BenchmarkGraphServer_DeleteGraphPerf(b *testing.B) {

	ctx := context.Background()
	s := newServer()

	gs := []*pb.Graph{
		{Vertices: []int32{1, 2, 3, 4, 5},
			Edges: map[int32]*pb.Neighbors{
				1: {Neighbors: []int32{2, 3, 4, 5}},
				2: {Neighbors: []int32{1}},
				3: {Neighbors: []int32{1}},
				4: {Neighbors: []int32{1}},
				5: {Neighbors: []int32{1}},
			}},
		{Vertices: []int32{1, 2, 3, 4, 5},
			Edges: map[int32]*pb.Neighbors{
				1: {Neighbors: []int32{2}},
				2: {Neighbors: []int32{1, 3}},
				3: {Neighbors: []int32{2, 4}},
				4: {Neighbors: []int32{3, 5}},
				5: {Neighbors: []int32{4}},
			}},
		{Vertices: []int32{1, 2, 3, 4, 5, 6, 7},
			Edges: map[int32]*pb.Neighbors{
				1: {Neighbors: []int32{2, 3, 4, 5}},
				2: {Neighbors: []int32{1}},
				3: {Neighbors: []int32{1}},
				4: {Neighbors: []int32{1}},
				5: {Neighbors: []int32{1}},
			}},
	}

	for i, g := range gs {
		id, err := s.PostGraph(ctx, g)

		if err == nil {
			if id.Id != int32(i+1) {
				b.Error("response: expected", i+1, "received", id.Id)
			}
		}
	}

	for i := 1; i <= 4; i++ {
		rep, _ := s.DeleteGraph(ctx, &pb.GraphID{Id: int32(i)})

		if i < 4 {
			if rep == nil {
				b.Error("delete failed when it should succeed")
			}
		} else {
			if rep != nil {
				b.Error("delete succeeded when it should fail")
			}
		}

	}

}
