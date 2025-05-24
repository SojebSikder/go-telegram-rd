package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func ExtractResourceID(url string) string {
	parts := strings.Split(url, "/")
	resourceID := parts[len(parts)-1]

	if strings.Contains(resourceID, "?") {
		resourceID = strings.Split(resourceID, "?")[0]
	}

	if strings.Contains(resourceID, ".") {
		resourceID = strings.Split(resourceID, ".")[0]
	}

	if strings.Contains(resourceID, "_") {
		resourceID = strings.Split(resourceID, "_")[1]
	}

	return resourceID
}

// DownloadFile fetches a Freepik resource by its ID and downloads it
func DownloadFile(url string) {
	resourceID := ExtractResourceID(url)

	apiKey := os.Getenv("FREEPIK_API_KEY")

	if apiKey == "" {
		fmt.Println("Missing FREEPIK_API_KEY environment variable")
		return
	}

	apiURL := fmt.Sprintf("https://api.freepik.com/v1/resources/%s/download", resourceID)

	// Prepare request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("x-freepik-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API error: %d\n%s\n", resp.StatusCode, string(body))
		return
	}

	// Decode JSON response
	var result struct {
		Data struct {
			Filename string `json:"filename"`
			URL      string `json:"url"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Validate URL
	if result.Data.URL == "" {
		fmt.Println("API did not return a valid download URL")
		return
	}

	// Download file
	downloadResp, err := http.Get(result.Data.URL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		fmt.Printf("Download failed with status %d\n", downloadResp.StatusCode)
		return
	}

	// Save file
	outFile, err := os.Create(result.Data.Filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, downloadResp.Body)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}

	fmt.Println("File downloaded successfully:", result.Data.Filename)
}
