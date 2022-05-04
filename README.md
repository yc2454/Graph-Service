# Graph-Service

## Structure
This project is consist of two parts: 1. a graph data struture and shortest path algorithm; 2. RPC server and client for graph-related services. 

The first part is located in the `go-graph/` directory. This part of the code is referenced from a [blog post](https://medium.com/@rishabhmishra131/golang-dijkstra-algorithm-7bf2722ba0c8), with some small modifications from me. 

The second part is located in the `graph_service/` directory. The service is defined by the protobuf file `graph.proto`, which contains 3 RPC services: `PostGraph`, `ShortestPath`, and `DeleteGraph`. I further implemented the server and the client code, as well as a unit test, a functional test, and a performance test. I protected the server operation with `sync.Mutex` so that it can support concurrent clients. The client and server code are in their respective folder, and the test are located together with the server.

## Running the Service
To run the service from command lines, first head to the `graph_service` directory and run start running the server:
```
cd graph_service
go mod tidy
go run server/server.go
```
You will see this output:
```
2022/05/03 21:23:42 server listening at [::]:8080
```
And then, in a separate terminal, run:
```
go run client/client.go
```
This starts a basic test case for the service. Here is a sample output:
```
2022/05/03 21:27:33 Posting graph
2022/05/03 21:27:33 Got ID: 3
2022/05/03 21:27:33 Asking for the shortest path between 2 and 3
2022/05/03 21:27:33 Received reply for shortest path
The shortest path between 2 and 3 is: 2 1 3 
2022/05/03 21:27:33 Deleting graph 3
2022/05/03 21:27:33 Successfully deleted the graph
``` 
The client posts a graph, queries about a shortest path, and then deletes the graph. 

I have also attempted a concurrent test for one client in `client_concurrent/client_concurrent.go`, where the client sends multiple requests to the server concurrently. Afterwards, we log down various evaluations of the server. A sample result of running `go run client_concurrent/client_concurrent.go` is:
```
2022/05/04 16:07:49 Posting graph 0
2022/05/04 16:07:49 Posting graph 1
2022/05/04 16:07:49 Posting graph 2
2022/05/04 16:07:49 Posting graph 3
2022/05/04 16:07:49 Posting graph 4
2022/05/04 16:07:49 Posting graph 5
2022/05/04 16:07:49 Posting graph 6
2022/05/04 16:07:49 Posting graph 7
2022/05/04 16:07:49 Posting graph 8
2022/05/04 16:07:49 Posting graph 9
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
2022/05/04 16:07:49 Successfully deleted the graph
Total Requests:             1000 hits
Availability:               100.00 %
Elapsed time:               0.02 secs
Request rate:               59502.41 trans/sec
Successful requests:        1000
Failed requests:            0
Longest request:            16222.46 us
Shortest request:           1505.58 us
Average request:            11524.33 us
Request std variance:       3060.13 us
```

More test cases are in the three test files. To run these tests, type
```
cd server
go test -v
```
The user can also use the `-bench` option to see the performance of the service, for example:
```
go test -bench=PostGraphPerf
```

## Future Directions
1. When we are finding the shortest path from S to T, we can also find the shortest path from S to all other nodes in the path. If we cache this result, then we can speed up future operations. The cache can also have evictions depending on the frequency of visits (evict the less frequently visited results).
2. We can add more complexity to the tests. For example, we can have multiple clients making requests concurrently, or we can add more randomization to 
graph and path generation.  