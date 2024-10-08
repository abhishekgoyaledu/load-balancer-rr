package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	for points := 1; points <= 100; points++ {
		status, body, err := processGameData(points)
		if err != nil {
			fmt.Printf("Error for points %d: %v\n", points, err)
			continue
		}

		fmt.Printf("Response status for points %d: %s\n", points, status)
		fmt.Printf("Response body: %s\n", body)

		time.Sleep(2 * time.Second) // Optionally, add a delay
	}
}

// processGameData sends game data as a POST request and returns the response status and body.
func processGameData(points int) (string, string, error) {
	gameData := map[string]interface{}{
		"game":    "Mobile Legends",
		"gamerID": "GYUTDTE",
		"points":  points,
	}

	jsonData, err := json.Marshal(gameData)
	if err != nil {
		return "", "", fmt.Errorf("error marshaling to JSON: %w", err)
	}

	resp, err := http.Post("http://localhost:8082/create", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", "", fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("error reading response body: %w", err)
	}

	return resp.Status, string(body), nil
}
