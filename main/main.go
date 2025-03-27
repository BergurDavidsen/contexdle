package main

import (
	"fmt"
	"os"
	"sort"

	Graph "github.com/BergurDavidsen/contexdle/Graph"
	"github.com/BergurDavidsen/contexdle/cache"
	"github.com/BergurDavidsen/contexdle/utils"
	levenshtein "github.com/agnivade/levenshtein"
)

func LevenshteinDistance(w1, w2 string) int {
	return levenshtein.ComputeDistance(w1, w2)
}

func GetKeys[T comparable, V any](m map[T]V) []T {
	keys := make([]T, 0, len(m)) // Preallocate memory
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func main() {

	graph := Graph.NewGraph(50)

	var words map[string]string

	fmt.Println("Trying to load cache")
	cachedSimilarityScores, err := cache.LoadCache()
	if err != nil {
		fmt.Println("Error loading cache:", err)
		return
	}

	var similarityScores map[string][]*Graph.EdgeScore

	if cachedSimilarityScores != nil {
		fmt.Println("Using cached similarity scores.")
		similarityScores = cachedSimilarityScores
		// Print the cached results
		fmt.Printf("loaded %d scores from cache\n", len(cachedSimilarityScores))
		// graph.Populate(cachedSimilarityScores)
		// graph.Print()
		//utils.PrintSortedSimilarityScores(similarityScores)
	} else {
		fmt.Println("No cache found. Computing similarity scores.")
		err := utils.LoadFile("../data/dictionary.json", &words)
		if err != nil {
			fmt.Println(err)
		}

		filteredWords := make(map[string]string)

		for w, d := range words {
			if utils.FilterWord(w) {
				graph.AddVertex(w)
				filteredWords[w] = d
			}
		}

		keyWords := GetKeys(filteredWords)

		fmt.Printf("Computing scores for %d words", len(keyWords))
		similarityScores = utils.ComputeSimilarityScoresParallel(keyWords)

		fmt.Println("Sorting the scores")
		for key := range similarityScores {
			sort.Sort(Graph.ByScore(similarityScores[key]))
		}

		// Save the computed scores into the cache
		fmt.Println("Saving cache")
		if err := cache.SaveCache(similarityScores); err != nil {
			fmt.Println("Error saving cache:", err)
		}

		fmt.Println("Saved cache and done!")

		// graph.Populate(similarityScores)
		// graph.Print()
	}

	fmt.Println("Populating graph")
	graph.Populate(similarityScores)

	if len(os.Args) > 2 {
		word1 := os.Args[1]
		word2 := os.Args[2]
		path, cost := graph.Dijkstra(word1, word2)

		fmt.Printf("path from %s -> %s: ", word1, word2)
		for _, w := range path {

			fmt.Printf("%s, ", w)
		}
		fmt.Printf("(cost: %f)\n", cost)
	}
}
