package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/yc2454/Graph-Service/graph_service"
)

/*
#include <stdlib.h>
*/
import "C"

var (
	addr = flag.String("addr", "localhost:8080", "the address to connect to")
)

func constructGraph() *pb.Graph {
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

	return g
}

// type StatusCode uint8
//
// const (
// 	Success StatusCode = iota
// 	Failure            = iota
// )

type QueryResult struct {
	queryID int
	path    *pb.Path
	err     error
}

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

	// Set a random seed
	C.srandom(C.uint(0))

	// Settings
	numGraphs := 10
	numQueries := 1000
	showResult := false

	// Initialize the a set of random graphs, saving the gids
	gids := make([]*pb.GraphID, numGraphs)

	for i := 0; i < numGraphs; i++ {
		g := constructGraph()

		log.Printf("Posting graph %d", i)
		id, err := c.PostGraph(ctx, g)

		if err != nil {
			log.Fatalf("could not post graph: %v", err)
		}

		gids[i] = id
	}

	// Start to posting requests
	var compQueue = make(chan QueryResult, 32)

	starts := make([]time.Time, numQueries)
	durations := make([]time.Duration, numQueries)
	start := time.Now()

	for i := 0; i < numQueries; i++ {
		// choose a random id from gids
		starts[i] = time.Now()
		idx := int(C.random()) % numGraphs
		go func(queryID int, gid *pb.GraphID) {
			req := new(pb.PathRequest)
			req.Gid = gid
			req.S = 2
			req.T = 3

			path, err := c.ShortestPath(ctx, req)

			if err != nil {
				log.Fatalf("could find shortest graph (id = %d): %v", gid, err)
			}

			compQueue <- QueryResult{queryID, path, err}
		}(i, gids[idx])
	}

	// Handle completions of tasks
	for i := 0; i < numQueries; i++ {
		res := <-compQueue
		if res.err != nil {
			log.Printf("could find shortest graph: queryID: %v, err: %v", res.queryID, err)
			continue
		}

		durations[res.queryID] = time.Now().Sub(starts[res.queryID])

		if showResult {
			fmt.Printf("queryID: %v, path: ", res.queryID)
			for _, n := range res.path.Path {
				fmt.Printf("%v ", n)
			}
			fmt.Println()
		}
	}

	end := time.Now()

	// Delete all graphs
	for _, id := range gids {
		reply, err := c.DeleteGraph(ctx, id)
		if err != nil {
			log.Fatalf("could find delete graph: %v", err)
		}

		log.Printf(reply.Result)
	}

	printStatistics(start, end, durations)
}

func printStatistics(start time.Time, end time.Time, latencies []time.Duration) {
	totalReqs := len(latencies)

	// Successful and failed requests
	numSuccesses := 0
	for _, d := range latencies {
		if d > 0 {
			numSuccesses += 1
		}
	}

	// Calc statistics
	avg := 0.0
	max := 0.0
	min := 1e9
	for _, lat := range latencies {
		if lat == 0 {
			continue
		}
		d := float64(lat)
		avg += d
		if max < d {
			max = d
		}
		if min > d {
			min = d
		}
	}

	avg /= float64(numSuccesses)

	// Calc std variance
	variance := 0.0
	for _, lat := range latencies {
		d := float64(lat)
		variance += (avg - d) * (avg - d)
	}

	variance /= float64(numSuccesses)
	variance = math.Sqrt(variance)

	fmt.Printf("Total Requests:             %d hits\n", totalReqs)
	fmt.Printf("Availability:               %.2f %%\n", 100*float64(numSuccesses)/float64(totalReqs))
	fmt.Printf("Elapsed time:               %.2f secs\n", end.Sub(start).Seconds())
	fmt.Printf("Request rate:               %.2f trans/sec\n", float64(totalReqs)/end.Sub(start).Seconds())
	fmt.Printf("Successful requests:        %d\n", numSuccesses)
	fmt.Printf("Failed requests:            %d\n", totalReqs-numSuccesses)
	fmt.Printf("Longest request:            %.2f us\n", max/1e3)
	fmt.Printf("Shortest request:           %.2f us\n", min/1e3)
	fmt.Printf("Average request:            %.2f us\n", avg/1e3)
	fmt.Printf("Request std variance:       %.2f us\n", variance/1e3)
}
