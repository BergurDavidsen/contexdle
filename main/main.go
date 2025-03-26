package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	graph "github.com/BergurDavidsen/contexdle/Graph"
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

type EdgeScore struct {
	Word  string
	Score float64
}

type ByScore []*EdgeScore

func (a ByScore) Len() int      { return len(a) }
func (a ByScore) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Sorting in descending order
func (a ByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }

var cacheFile = "similarity_cache.gob"

// Save the cache to a file using gob
func saveCache(cache map[string][]*EdgeScore) error {
	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new encoder and write the cache to the file
	encoder := gob.NewEncoder(file)
	return encoder.Encode(cache)
}

func loadCache() (map[string][]*EdgeScore, error) {
	// Check if the cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, nil // Cache doesn't exist, return nil
	}

	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cache map[string][]*EdgeScore
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, err
	}

	return cache, nil
}

// ComputeLevenshteinSimilarity calculates similarity based on Levenshtein distance
func ComputeLevenshteinSimilarity(word1, word2 string) float64 {
	levDistance := LevenshteinDistance(word1, word2) // Assume you have this function implemented
	maxLen := max(len(word1), len(word2))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(levDistance)/float64(maxLen)
}

// ComputeSimilarityScores finds similar words and stores meaningful relationships
func ComputeSimilarityScoresParallel(words []string) map[string][]*EdgeScore {
	similarityScores := make(map[string][]*EdgeScore)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var completed int64 // Atomic counter

	totalComparisons := len(words) * (len(words) - 1) / 2 // Total word pairs

	for i, word := range words {
		wg.Add(1)
		go func(i int, word string) {
			defer wg.Done()
			for j := i + 1; j < len(words); j++ {
				other := words[j]
				score := ComputeLevenshteinSimilarity(word, other)

				if score > 0.3 {
					mu.Lock()
					similarityScores[word] = append(similarityScores[word], &EdgeScore{Word: other, Score: score})
					similarityScores[other] = append(similarityScores[other], &EdgeScore{Word: word, Score: score})
					mu.Unlock()
				}

				// Track progress
				current := atomic.AddInt64(&completed, 1)
				if current%1000 == 0 { // Print every 1000 comparisons
					fmt.Printf("\rProgress: %.2f%%", (float64(current)/float64(totalComparisons))*100)
				}
			}
		}(i, word)
	}
	wg.Wait()
	fmt.Println("\nProcessing complete!")
	return similarityScores
}

func FilterWord(w string) bool {
	if len(w) < 3 || len(w) > 8 {
		return false
	}

	// Convert the word to lowercase for consistent comparison
	w = strings.ToLower(w)

	// Check if the word contains spaces or hyphens
	if strings.Contains(w, " ") || strings.Contains(w, "-") {
		return false
	}

	// Check if the word contains only alphabetic characters (no numbers or special chars)
	for _, r := range w {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func loadFile(filename string, structure any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file")
		return err
	}
	return json.Unmarshal(data, structure)
}
func main() {

	graph := graph.NewGraph()

	var words map[string]string

	fmt.Println("Trying to load cache")
	cachedSimilarityScores, err := loadCache()
	if err != nil {
		fmt.Println("Error loading cache:", err)
		return
	}

	if cachedSimilarityScores != nil {
		fmt.Println("Using cached similarity scores.")
		similarityScores := cachedSimilarityScores

		// Print the cached results
		fmt.Printf("loaded %d scores from cache", len(similarityScores))
	} else {
		fmt.Println("No cache found. Computing similarity scores.")
		err := loadFile("../data/dictionary.json", &words)
		if err != nil {
			fmt.Println(err)
		}

		filteredWords := make(map[string]string)

		for w, d := range words {
			if FilterWord(w) {
				graph.AddVertex(w, d)
				filteredWords[w] = d
			}
		}

		keyWords := GetKeys(filteredWords)

		fmt.Printf("Computing scores for %d words", len(keyWords))
		similarityScores := ComputeSimilarityScoresParallel(keyWords)

		fmt.Println("Sorting the scores")
		for key := range similarityScores {
			sort.Sort(ByScore(similarityScores[key]))
		}

		// Save the computed scores into the cache
		fmt.Println("Saving cache")
		if err := saveCache(similarityScores); err != nil {
			fmt.Println("Error saving cache:", err)
		}

		fmt.Println("Saved cache and done!")
	}

}

func printSortedSimilarityScores(similarityScores map[string][]*EdgeScore) {
	for word, edges := range similarityScores {
		fmt.Printf("%s -> ", word)
		for i, edge := range edges {
			if 15 <= i {
				break
			}
			fmt.Printf("(%s, %.2f) ", edge.Word, edge.Score)
		}
		fmt.Println()
	}
}
