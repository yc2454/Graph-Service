package graphservice

import (
	"context"
	"errors"
	"flag"

	graph "github.com/yc2454/Graph-Service/graph"

	pb "github.com/yc2454/Graph-Service/graph_service"
)

var (
	port = flag.Int("port", 8080, "The server port")
)

type graphServiceServer struct {
	pb.UnimplementedGraphServiceServer
	graphs map[int32]*graph.ItemGraph
	curID  int32

	// mu         sync.Mutex // protects routeNotes
	// routeNotes map[string][]*pb.RouteNote
}

func (s *graphServiceServer) PostGraph(ctx context.Context, g *pb.Graph) (*pb.GraphID, error) {
	newGraph := graph.NewGraph()

	for _, v := range g.GetVertices() {
		n := graph.NewNode(int(v))
		newGraph.AddNode(n)
	}

	edges := g.GetEdges()

	// newGraph.String()

	for _, v := range g.GetVertices() {
		if edges[v] != nil {
			for _, u := range edges[v].Neighbors {

				n1, err1 := newGraph.FindNode(int(v))
				n2, err2 := newGraph.FindNode(int(u))

				if err1 != nil || err2 != nil {
					return nil, errors.New("found edge between non-existant nodes")
				} else {
					newGraph.AddEdge(n1, n2)
				}

			}
		}
	}

	s.graphs[s.curID] = newGraph
	id := new(pb.GraphID)
	id.Id = int32(s.curID)
	s.curID++

	return id, nil
}

func (s *graphServiceServer) ShortestPath(ctx context.Context, req *pb.PathRequest) (*pb.Path, error) {

	g := s.graphs[req.Gid.Id]

	// g.String()

	n1, err1 := g.FindNode(int(req.S))
	n2, err2 := g.FindNode(int(req.T))

	res := new(pb.Path)

	if err1 == nil && err2 == nil {
		p, _ := g.GetShortestPath(n1, n2)

		for _, n := range p {
			res.Path = append(res.Path, int32(n))
		}
		return res, nil

	} else {
		return nil, errors.New("non-existant node")
	}
}

func (s *graphServiceServer) DeleteGraph(ctx context.Context, id *pb.GraphID) (*pb.DeleteReply, error) {

	g := s.graphs[id.Id]
	reply := new(pb.DeleteReply)

	if g == nil {
		reply.Result = "No such id stored"
	} else {
		s.graphs[id.Id] = nil
		reply.Result = "Successfully deleted the graph"
	}

	return reply, nil
}

func newServer() *graphServiceServer {
	s := new(graphServiceServer)
	s.graphs = make(map[int32]*graph.ItemGraph)
	s.curID = 1
	return s
}

// func main() {
// 	flag.Parse()
// 	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
// 	if err != nil {
// 		log.Fatalf("failed to listen: %v", err)
// 	}

// 	s := grpc.NewServer()
// 	pb.RegisterGraphServiceServer(s, newServer())

// 	log.Printf("server listening at %v", lis.Addr())
// 	if err := s.Serve(lis); err != nil {
// 		log.Fatalf("failed to serve: %v", err)
// 	}
// }
