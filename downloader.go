package main

import (
	"encoding/json"
	"errors"
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
func DownloadFile(url string, isBytes bool) ([]byte, string, error) {
	resourceID := ExtractResourceID(url)

	apiKey := os.Getenv("FREEPIK_API_KEY")

	if apiKey == "" {
		fmt.Println("Missing FREEPIK_API_KEY environment variable")
		return nil, "", errors.New("missing FREEPIK_API_KEY environment variable")
	}

	apiURL := fmt.Sprintf("https://api.freepik.com/v1/resources/%s/download", resourceID)

	// Prepare request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, "", errors.New("error creating request")
	}
	req.Header.Set("x-freepik-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error executing request:", err)
		return nil, "", errors.New("error executing request")
	}
	defer resp.Body.Close()

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("API error: %d\n%s\n", resp.StatusCode, string(body))
		return nil, "", errors.New("API error")
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
		return nil, "", errors.New("error decoding JSON")
	}

	// Validate URL
	if result.Data.URL == "" {
		fmt.Println("API did not return a valid download URL")
		return nil, "", errors.New("API did not return a valid download URL")
	}

	// Download file
	downloadResp, err := http.Get(result.Data.URL)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return nil, "", errors.New("error downloading file")
	}
	defer downloadResp.Body.Close()

	if downloadResp.StatusCode != http.StatusOK {
		fmt.Printf("Download failed with status %d\n", downloadResp.StatusCode)
		return nil, "", errors.New("download failed")
	}

	// if isBytes is true, return the file as bytes
	if isBytes {
		fileBytes, _ := io.ReadAll(downloadResp.Body)
		return fileBytes, result.Data.Filename, nil
	}

	// Save file
	outFile, err := os.Create(result.Data.Filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, "", errors.New("error creating file")
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, downloadResp.Body)
	if err != nil {
		fmt.Println("Error saving file:", err)
		return nil, "", errors.New("error saving file")
	}

	return nil, "", nil
}
