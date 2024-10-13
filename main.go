package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Data struct {
	A int `json:"a"`
	B int `json:"b"`
}

func readFile(filepath string) ([]Data, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var data []Data
	if err := json.Unmarshal(file, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func sumData(data []Data, start int, end int, sumChan chan int) {
	sum := 0
	for i := start; i < end; i++ {
		sum += data[i].A + data[i].B
	}
	sumChan <- sum
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <num_goroutines>\n", os.Args[0])
	}

	numGoroutines, err := strconv.Atoi(os.Args[1])
	if err != nil || numGoroutines <= 0 {
		log.Fatalf("Invalid number of goroutines: %s\n", os.Args[1])
	}

	// Get the current directory
	projectRoot, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v\n", err)
	}

	dataFilePath := filepath.Join(projectRoot, "data", "data.json")

	data, err := readFile(dataFilePath)
	if err != nil {
		log.Fatalf("Error reading file: %v\n", err)
	}

	resultChan := make(chan int, numGoroutines)
	blockSize := len(data) / numGoroutines

	for i := 0; i < numGoroutines; i++ {
		start := i * blockSize
		end := start + blockSize
		if i == numGoroutines-1 {
			end = len(data)
		}
		go sumData(data, start, end, resultChan)
	}

	totalSum := 0
	for i := 0; i < numGoroutines; i++ {
		totalSum += <-resultChan
	}

	fmt.Printf("Total Sum: %d\n", totalSum)
}
