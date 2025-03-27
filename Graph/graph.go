package graph

import (
	"container/heap"
	"fmt"
	"math"
)

type Vertex struct {
	Val string
}

type Edge struct {
	Weight float32
	From   *Vertex
	To     *Vertex
}

type EdgeScore struct {
	Word  string
	Score float32
}

type ByScore []*EdgeScore

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

type Graph struct {
	Vertices   map[string]*Vertex
	Edges      map[string][]*Edge
	MaxDegrees int
}

// Priority Queue for Dijkstra
type Item struct {
	Vertex   string
	Distance float32
	Index    int
}

type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Distance < pq[j].Distance // Min-heap based on distance
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func (graph *Graph) AddVertex(value string) {
	graph.Vertices[value] = &Vertex{Val: value}
	graph.Edges[value] = []*Edge{}
}

func (graph *Graph) AddEdge(from, to string, weight float32) {
	// check if src & dest exist
	if _, ok := graph.Vertices[from]; !ok {
		return
	}
	if _, ok := graph.Vertices[to]; !ok {
		return
	}
	// add edge src --> dest
	edge := &Edge{Weight: weight, From: graph.Vertices[from], To: graph.Vertices[to]}
	graph.Edges[from] = append(graph.Edges[from], edge)
}

func (graph *Graph) GetNeighbors(value string) []*Edge {
	// check if Vertex exists
	neighbors := graph.Edges[value]
	return neighbors
}

func NewGraph(maxDegrees int) *Graph {
	g := &Graph{Vertices: map[string]*Vertex{}, Edges: map[string][]*Edge{}, MaxDegrees: maxDegrees}
	return g
}

func (graph *Graph) Populate(similarityScores map[string][]*EdgeScore) {
	for word, _ := range similarityScores {
		graph.AddVertex(word)
	}
	for word, edges := range similarityScores {
		for i, edge := range edges {
			if graph.MaxDegrees <= i {
				break
			}
			graph.AddEdge(word, edge.Word, edge.Score)
		}
	}
}

// Dijkstra's algorithm
func (graph *Graph) Dijkstra(startWord, endWord string) ([]string, float32, error) {
	// Initialize data structures
	distances := make(map[string]float32)
	previous := make(map[string]string)
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Check if both values exist in graph
	if _, ok := graph.Vertices[startWord]; !ok {
		return nil, 0, fmt.Errorf("'%s' does not exist in the graph", startWord)
	}
	if _, ok := graph.Vertices[endWord]; !ok {
		return nil, 0, fmt.Errorf("'%s' does not exist in the graph", endWord)
	}
	// Set all distances to infinity
	for vertex := range graph.Vertices {
		distances[vertex] = float32(math.Inf(1))
		previous[vertex] = ""
	}

	// Set the start vertex's distance to 0
	distances[startWord] = 0
	heap.Push(&pq, &Item{
		Vertex:   startWord,
		Distance: 0,
	})

	// Dijkstra's algorithm loop
	for pq.Len() > 0 {
		// Get the vertex with the smallest distance
		current := heap.Pop(&pq).(*Item)
		currentVertex := current.Vertex

		// If we reached the destination, we can stop
		if currentVertex == endWord {
			break
		}

		// Visit each neighbor of the current vertex
		for _, edge := range graph.GetNeighbors(currentVertex) {
			neighbor := edge.To.Val
			newDist := distances[currentVertex] + float32(edge.Weight)

			// If a shorter path to the neighbor is found
			if newDist < distances[neighbor] {
				distances[neighbor] = newDist
				previous[neighbor] = currentVertex
				heap.Push(&pq, &Item{
					Vertex:   neighbor,
					Distance: newDist,
				})
			}
		}
	}

	// Reconstruct the shortest path
	path := []string{}
	for currentVertex := endWord; currentVertex != ""; currentVertex = previous[currentVertex] {
		path = append([]string{currentVertex}, path...)
	}

	// If there's no path to the endWord, return an empty path
	if len(path) == 1 && path[0] != endWord {
		return nil, float32(math.Inf(1)), nil
	}

	return path, distances[endWord], nil
}

func (graph *Graph) Print() {
	fmt.Println("The Graph")
	fmt.Println("---------")
	for v := range graph.Vertices {
		fmt.Printf("%s -> ", v)

		neighbors := graph.Edges[v]

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
