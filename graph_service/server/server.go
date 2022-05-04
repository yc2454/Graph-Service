package graphservice

import (
	"context"
	"errors"
	"sync"

	graph "github.com/yc2454/Graph-Service/graph"

	pb "github.com/yc2454/Graph-Service/graph_service"
)

type graphServiceServer struct {
	pb.UnimplementedGraphServiceServer

	// A mapping from graph ID to stored graphs
	graphs map[int32]*graph.ItemGraph

	// The next ID to return for newly posted graph.
	// Monotonically increase from 1
	curID int32

	mu sync.Mutex // protects routeNotes
}

// PostGraph receives a graph from the client, stores it in the server with an
// ID, and return the graph ID if the graph is valid.
func (s *graphServiceServer) PostGraph(ctx context.Context, g *pb.Graph) (*pb.GraphID, error) {

	// Initialize the graph to store
	newGraph := graph.NewGraph()

	// Record the nodes in the graph to post
	for _, v := range g.GetVertices() {
		n := graph.NewNode(int(v))
		newGraph.AddNode(n)
	}

	edges := g.GetEdges()

	// Connect the edges in the graph to post
	for _, v := range g.GetVertices() {
		if edges[v] != nil {
			for _, u := range edges[v].Neighbors {

				// First, retrieve the nodes from the graph
				n1, err1 := newGraph.FindNode(int(v))
				n2, err2 := newGraph.FindNode(int(u))

				if err1 != nil || err2 != nil {
					return nil, errors.New("found edge between non-existant nodes")
				} else {
					// Connect the edge if both nodes have been recorded
					newGraph.AddEdge(n1, n2)
				}

			}
		}
	}

	s.mu.Lock()
	s.graphs[s.curID] = newGraph
	id := new(pb.GraphID)
	id.Id = int32(s.curID)

	// Increase the ID
	s.curID++
	s.mu.Unlock()

	return id, nil
}

// ShortestPath takes the path request from the client, which contains a graph ID
// and the start and end point of the path. It returns the shortest path if such a
// path exists.
func (s *graphServiceServer) ShortestPath(ctx context.Context, req *pb.PathRequest) (*pb.Path, error) {

	g := s.graphs[req.Gid.Id]

	// Retrieve the nodes from the graph first
	n1, err1 := g.FindNode(int(req.S))
	n2, err2 := g.FindNode(int(req.T))

	// [res] is used to store the result to return
	res := new(pb.Path)

	if err1 == nil && err2 == nil {
		// Compute the shortest path
		p, _ := g.GetShortestPath(n1, n2)

		// Record the path in [res]
		for _, n := range p {
			res.Path = append(res.Path, int32(n))
		}
		return res, nil

	} else {
		// When we cannot retrieve the nodes, return error
		return nil, errors.New("non-existant node")
	}
}

// DeleteGraph deletes the graph with ID=[id] from the server and
// returns a message to the client if such graph exists.
func (s *graphServiceServer) DeleteGraph(ctx context.Context, id *pb.GraphID) (*pb.DeleteReply, error) {

	g := s.graphs[id.Id]
	reply := new(pb.DeleteReply)

	if g == nil {
		return nil, errors.New("non-existant graph")
	} else {
		s.mu.Lock()
		s.graphs[id.Id] = nil
		s.mu.Unlock()
		reply.Result = "Successfully deleted the graph"
		return reply, nil
	}

}

// Constructor of the server
func newServer() *graphServiceServer {
	s := new(graphServiceServer)
	s.graphs = make(map[int32]*graph.ItemGraph)
	s.curID = 1
	return s
}
