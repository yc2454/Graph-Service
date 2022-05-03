module github.com/yc2454/Graph-Service/graph_service

go 1.18

require (
	github.com/yc2454/Graph-Service/graph v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.46.0
	google.golang.org/protobuf v1.28.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/yc2454/Graph-Service v0.0.0-20220503194849-00e781a0fd67 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.3 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)

replace github.com/yc2454/Graph-Service/graph => ../go_graph
