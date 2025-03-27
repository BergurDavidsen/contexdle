package cache

import (
	"encoding/gob"
	"os"

	graph "github.com/BergurDavidsen/contexdle/Graph"
)

var cacheFile = "cache/default_similarity_cache.gob"

func SaveCache(cache map[string][]*graph.EdgeScore) error {
	file, err := os.Create(cacheFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a new encoder and write the cache to the file
	encoder := gob.NewEncoder(file)
	return encoder.Encode(cache)
}

func LoadCache() (map[string][]*graph.EdgeScore, error) {
	// Check if the cache file exists
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return nil, nil // Cache doesn't exist, return nil
	}

	file, err := os.Open(cacheFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cache map[string][]*graph.EdgeScore
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&cache); err != nil {
		return nil, err
	}

	return cache, nil
}
