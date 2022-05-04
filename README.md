# Graph-Service

## Structure
This project is consist of two parts: 1. a graph data struture and shortest path algorithm; 2. RPC server and client for graph-related services. 

The first part is located in the `go-graph/` directory. This part of the code is referenced from a [blog post](https://medium.com/@rishabhmishra131/golang-dijkstra-algorithm-7bf2722ba0c8), with some small modifications from me. 

The second part is located in the `graph_service/` directory. The service is defined by the protobuf file `graph.proto`, which contains 3 RPC services: `PostGraph`, `ShortestPath`, and `DeleteGraph`. I further implemented the server and the client code, as well as a unit test, a functional test, and a performance test. The client and server code are in their respective folder, and the test are located together with the server.

## Running the Service
To run the service from command lines, first head to the `graph_service` directory and run start running the server:
```
cd graph_service
go run server/server.go
```
You will see this output:
```
2022/05/03 21:23:42 server listening at [::]:8080
```
And then, in a separate terminal, run:
```
go run client/client
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

More test cases are in the three test files. To run these tests, type
```
go test
```
The user can also use the `-bench` option to see the performance of the service, for example:
```
go test -bench=PostGraphPerf
```