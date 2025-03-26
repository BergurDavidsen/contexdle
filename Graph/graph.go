package graph

import "fmt"

type Vertex struct {
	Val string
	Def string
}
type Edge struct {
	Weight int
	From   *Vertex
	To     *Vertex
}

type Graph struct {
	Vertices map[string]*Vertex
	Edges    map[string][]*Edge
}

func (graph *Graph) AddVertex(value, definition string) {
	graph.Vertices[value] = &Vertex{Val: value}
	graph.Edges[value] = []*Edge{}
}

func (graph *Graph) AddEdge(from, to string, weight int) {
	// check if src & dest exist
	if _, ok := graph.Vertices[from]; !ok {
		return
	}
	if _, ok := graph.Vertices[to]; !ok {
		return
	}
	// add edge src --> dest
	edge := &Edge{Weight: weight, From: graph.Vertices[from], To: graph.Vertices[to]}
	graph.Edges[from] = append(graph.Edges[to], edge)
}

func (graph *Graph) GetNeighbors(value string) []*Edge {
	// check if Vertex exists
	neighbors := graph.Edges[value]
	return neighbors

}

func NewGraph() *Graph {
	g := &Graph{Vertices: map[string]*Vertex{}, Edges: map[string][]*Edge{}}

	return g
}

func (graph *Graph) PrintGraph() {
	fmt.Println("The Graph")
	fmt.Println("---------")
	for v := range graph.Vertices {
		fmt.Printf("%s -> ", v)

		neighbors := graph.GetNeighbors(v)

		if len(neighbors) == 0 {
			fmt.Println("[]")
			continue
		}

		for _, n := range neighbors {
			fmt.Printf("%s, ", n.To.Val)
		}
		fmt.Println()
	}
}
