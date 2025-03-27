package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	graph "github.com/BergurDavidsen/contexdle/Graph"
	"github.com/agnivade/levenshtein"
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

// ComputeLevenshteinSimilarity calculates similarity based on Levenshtein distance
func ComputeLevenshteinSimilarity(word1, word2 string) float32 {
	levDistance := LevenshteinDistance(word1, word2) // Assume you have this function implemented
	maxLen := max(len(word1), len(word2))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float32(levDistance)/float32(maxLen)
}

// ComputeSimilarityScores finds similar words and stores meaningful relationships
func ComputeSimilarityScoresParallel(words []string) map[string][]*graph.EdgeScore {
	similarityScores := make(map[string][]*graph.EdgeScore)
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
					similarityScores[word] = append(similarityScores[word], &graph.EdgeScore{Word: other, Score: score})
					similarityScores[other] = append(similarityScores[other], &graph.EdgeScore{Word: word, Score: score})
					mu.Unlock()
				}

				// Track progress
				current := atomic.AddInt64(&completed, 1)
				if current%1000 == 0 { // Print every 1000 comparisons
					fmt.Printf("\rProgress: %.2f%%", (float32(current)/float32(totalComparisons))*100)
				}
			}
		}(i, word)
	}
	wg.Wait()
	fmt.Println("\nProcessing complete!")
	return similarityScores
}

func FilterWord(w string) bool {
	// filters the length. The longer the more words will be added and more memory needed
	if len(w) < 3 || len(w) > 4 {
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

func LoadFile(filename string, structure any) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Could not read file")
		return err
	}
	return json.Unmarshal(data, structure)
}

func PrintSortedSimilarityScores(similarityScores map[string][]*graph.EdgeScore) {
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
