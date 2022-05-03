// Package graph creates a ItemGraph data structure for the Item type
package graph

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

// Node is the node in the graph.
// It is accompanied with an int value
type Node struct {
	value int
}

func (n *Node) String() string {
	return fmt.Sprintf("%v", n.value)
}

func (n *Node) Value() int {
	return n.value
}

// ItemGraph is the Items graph
type ItemGraph struct {
	nodes []*Node
	edges map[Node][]*Node
	lock  sync.RWMutex
}

// Constructor of the graph
func NewGraph() *ItemGraph {
	g := new(ItemGraph)
	return g
}

func NewNode(v int) *Node {
	n := new(Node)
	n.value = v
	return n
}

func (g *ItemGraph) FindNode(v int) (*Node, error) {
	for _, n := range g.nodes {
		if n.value == v {
			return n, nil
		}
	}
	return nil, errors.New("no such node")
}

// AddNode adds a node to the graph
func (g *ItemGraph) AddNode(n *Node) {
	g.lock.Lock()
	g.nodes = append(g.nodes, n)
	g.lock.Unlock()
}

// AddEdge adds an edge to the graph
func (g *ItemGraph) AddEdge(n1, n2 *Node) {
	g.lock.Lock()
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*n1] = append(g.edges[*n1], n2)
	g.edges[*n2] = append(g.edges[*n2], n1)
	g.edges[*n1] = unique(g.edges[*n1])
	g.edges[*n2] = unique(g.edges[*n2])
	g.lock.Unlock()
}

// Remove duplicate elements
func unique(intSlice []*Node) []*Node {
	keys := make(map[*Node]bool)
	list := []*Node{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// Print out the graph
func (g *ItemGraph) String() {
	g.lock.RLock()
	s := ""
	for i := 0; i < len(g.nodes); i++ {
		s += g.nodes[i].String() + " -> "
		near := g.edges[*g.nodes[i]]
		for j := 0; j < len(near); j++ {
			s += near[j].String() + " "
		}
		s += "\n"
	}
	fmt.Println(s)
	g.lock.RUnlock()
}

type Array []*Node

func (arr Array) hasPropertyOf(n *Node) bool {
	for _, v := range arr {
		if n.value == v.value {
			return true
		}
	}
	return false
}

func (graph *ItemGraph) ShortestPath(start *Node, end *Node, path Array) []*Node {

	if _, exist := graph.edges[*start]; !exist {
		return path
	}

	path = append(path, start)
	if start == end {
		return path
	}

	shortest := make([]*Node, 0)

	for _, node := range graph.edges[*start] {
		if !path.hasPropertyOf(node) {
			newPath := graph.ShortestPath(node, end, path)
			if len(newPath) > 0 {
				if len(shortest) == 0 || (len(newPath) < len(shortest)) {
					shortest = newPath
				}
			}
		}
	}

	return shortest
}

type Vertex struct {
	Node     *Node
	Distance int
}

func (g *ItemGraph) GetShortestPath(startNode *Node, endNode *Node) ([]int, int) {
	visited := make(map[int]bool)
	dist := make(map[int]int)
	prev := make(map[int]int)

	q := NodeQueue{}
	pq := q.NewQ()
	start := Vertex{
		Node:     startNode,
		Distance: 0,
	}
	for _, nval := range g.nodes {
		dist[nval.Value()] = math.MaxInt64
	}
	dist[startNode.Value()] = start.Distance
	pq.Enqueue(start)

	for !pq.IsEmpty() {
		v := pq.Dequeue()
		if visited[v.Node.Value()] {
			continue
		}
		visited[v.Node.Value()] = true
		near := g.edges[*v.Node]

		for _, val := range near {
			if !visited[val.Value()] {
				if dist[v.Node.Value()]+1 < dist[val.Value()] {
					store := Vertex{
						Node:     val,
						Distance: dist[v.Node.Value()] + 1,
					}
					dist[val.Value()] = dist[v.Node.Value()] + 1
					//prev[val.Node.Value] = fmt.Sprintf("->%s", v.Node.Value)
					prev[val.Value()] = v.Node.Value()
					pq.Enqueue(store)
				}
				//visited[val.Node.value] = true
			}
		}
	}

	// fmt.Println(dist)
	// fmt.Println(prev)
	pathval := prev[endNode.Value()]

	var finalArr []int
	finalArr = append(finalArr, endNode.Value())
	for pathval != startNode.Value() {
		// fmt.Printf("looping")
		finalArr = append(finalArr, pathval)
		pathval = prev[pathval]
	}
	finalArr = append(finalArr, pathval)
	// fmt.Println(finalArr)
	for i, j := 0, len(finalArr)-1; i < j; i, j = i+1, j-1 {
		finalArr[i], finalArr[j] = finalArr[j], finalArr[i]
	}
	return finalArr, dist[endNode.Value()]

}
