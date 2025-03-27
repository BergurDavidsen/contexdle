# WordGraph - A Graph-Based Word Similarity Finder

## Overview

**WordGraph** is a Go-based project that constructs a graph of words and connects them using a similarity score. The project primarily leverages **Levenshtein distance** to determine similarity between words and then applies **Dijkstra's algorithm** to find the shortest path between two words in the graph.  

This approach allows for an exploration of relationships between words based on their edit distance, providing a way to determine how "close" two words are in terms of their structure.  

## Features

- **Graph-Based Representation:** Words are represented as nodes, and edges represent their similarity scores.  
- **Levenshtein Distance Calculation:** Measures the edit distance between two words to establish connections.  
- **Dijkstra’s Algorithm for Shortest Path:** Finds the most efficient word transformation sequence.  
- **Parallelized Computation:** Uses Go’s concurrency features to speed up similarity score calculations.  
- **Caching Mechanism:** Saves computed similarity scores to avoid redundant calculations on repeated runs.  
- **Filtering Mechanism:** Ensures only valid words (length 2-7, no spaces, hyphens, or non-alphabet characters) are considered.  

## Installation & Usage

### Prerequisites

- **Go 1.18+**  
- A JSON dictionary file (e.g., `dictionary.json`) containing words  

### Installation

1. Clone the repository:  

   ```sh
   git clone https://github.com/BergurDavidsen/contexdle.git
   cd contexdle
   ```  

2. Install dependencies (if any).  

### Running the Program

To build and execute the program, use:  

```sh
go run main.go [word1] [word2]
```  

Example:  

```sh
go run main.go cat dog
```  

This will output the shortest path (if any) from "cat" to "dog" based on similarity scores.  

### Caching

The program attempts to load cached similarity scores before computing new ones. If no cache is found, it processes words from `dictionary.json`, computes similarity scores, and saves them for future use.  

## Project Structure

```
/contexdle
│── /Graph                # Graph implementation (vertices, edges, Dijkstra’s algorithm)
│── /utils                # Utility functions (similarity calculation, word filtering, etc.)
│── /cache                # Caching logic for similarity scores
│── main.go               # Entry point of the program
│── README.md             # Project documentation
│── dictionary.json       # Source data for words (not included in repo)
```

## Future Enhancements

- **Word2Vec Integration:** Improve similarity scoring by incorporating context-based embeddings.  
- **Hybrid Similarity Measure:** Combine Levenshtein distance with cosine similarity from Word2Vec.  
- **Web API or UI:** Build a simple interface for querying word relationships dynamically.  

## License

This project is open-source under the **MIT License**. Contributions are welcome!  
