package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type FilesRequestBody struct {
	StartEvent InputFileEvent `json:"start_event"`
	Context    map[string]any `json:"context"`
	HandlerId  string         `json:"handler_id"`
}

type InputFileEvent struct {
	FileId   string `json:"file_id"`
	Username string `json:"username"`
	FileName string `json:"file_name"`
}

type FilesResponseResult struct {
	Success bool    `json:"success"`
	Error   *string `json:"error"`
}

type FilesResponseBody struct {
	HandlerId    string               `json:"handler_id"`
	WorkflowName string               `json:"workflow_name"`
	RunId        string               `json:"run_id"`
	Status       string               `json:"status"`
	StartedAt    *string              `json:"started_at"`
	UpdatedAt    *string              `json:"updated_at"`
	CompletedAt  *string              `json:"completed_at"`
	Error        *string              `json:"error"`
	Result       *FilesResponseResult `json:"result"`
}

type SearchRequestBody struct {
	StartEvent SearchInputEvent `json:"start_event"`
	Context    map[string]any   `json:"context"`
	HandlerId  string           `json:"handler_id"`
}

type SearchInputEvent struct {
	SearchType  string  `json:"search_type"`
	SearchInput string  `json:"search_input"`
	Username    string  `json:"username"`
	FileName    *string `json:"file_name"`
	Category    *string `json:"category"`
}

type SearchResult struct {
	ResultType string  `json:"result_type"`
	Text       string  `json:"text"`
	Similarity float64 `json:"similarity"`
	FileName   string  `json:"file_name"`
	Category   string  `json:"category"`
}

type SearchResponseResult struct {
	Results []SearchResult `json:"results"`
}

type SearchResponseBody struct {
	HandlerId    string                `json:"handler_id"`
	WorkflowName string                `json:"workflow_name"`
	RunId        string                `json:"run_id"`
	Status       string                `json:"status"`
	StartedAt    *string               `json:"started_at"`
	UpdatedAt    *string               `json:"updated_at"`
	CompletedAt  *string               `json:"completed_at"`
	Error        *string               `json:"error"`
	Result       *SearchResponseResult `json:"result"`
}

func ProcessFile(fileInput InputFileEvent) (*FilesResponseBody, error) {
	requestBody := FilesRequestBody{StartEvent: fileInput, Context: map[string]any{}, HandlerId: ""}
	apiKey := os.Getenv("LLAMA_CLOUD_API_KEY")
	apiEndpoint := os.Getenv("FILES_API_ENDPOINT")
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	// Create the HTTP request
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response FilesResponseBody

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(*response.Error)
	}
	return &response, nil
}

func ProcessSearch(searchInput SearchInputEvent) (*SearchResponseBody, error) {
	requestBody := SearchRequestBody{StartEvent: searchInput, Context: map[string]any{}, HandlerId: ""}
	apiKey := os.Getenv("LLAMA_CLOUD_API_KEY")
	apiEndpoint := os.Getenv("SEARCH_API_ENDPOINT")
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	// Create the HTTP request
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "POST", apiEndpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response SearchResponseBody

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(*response.Error)
	}
	return &response, nil
}
