package agent

import (
	"os"
	"testing"

	"github.com/run-llama/study-llama/frontend/files"
)

func TestProcessFile(t *testing.T) {
	_, okApi := os.LookupEnv("LLAMA_CLOUD_API_KEY")
	user, okUser := os.LookupEnv("TEST_USER")
	_, okEndpoint := os.LookupEnv("FILES_API_ENDPOINT")
	if !okApi || !okUser || !okEndpoint {
		t.Skip("Necessary env variables not available")
	}
	file := "../testfiles/the-future-of-vibe-coding.pdf"
	src, _ := os.Open(file)
	fileId, err := files.UploadFile(src, file)
	if err != nil {
		t.Errorf("Expected no error while uploading the file, got %s", err.Error())
	}
	inputEvent := InputFileEvent{FileName: file, FileId: fileId, Username: user}
	res, err := ProcessFile(inputEvent)
	if err != nil {
		t.Errorf("Expected no error while processing the file, got %s", err.Error())
	}
	if res.GetErrorString() != nil {
		t.Errorf("Expected no error from the backend, got %s", *res.GetErrorString())
	}
}

func TestProcessSearch(t *testing.T) {
	_, okApi := os.LookupEnv("LLAMA_CLOUD_API_KEY")
	user, okUser := os.LookupEnv("TEST_USER")
	_, okEndpoint := os.LookupEnv("FILES_API_ENDPOINT")
	if !okApi || !okUser || !okEndpoint {
		t.Skip("Necessary env variables not available")
	}
	file := "../testfiles/the-future-of-vibe-coding.pdf"
	category := "vibecoding"
	inputEvent := SearchInputEvent{Username: user, FileName: &file, Category: &category, SearchType: "faqs", SearchInput: "What are the main risks associated with vibe-coding?"}
	res, err := ProcessSearch(inputEvent)
	if err != nil {
		t.Errorf("Expected no error while processing the file, got %s", err.Error())
	}
	if len(res.GetResults()) == 0 {
		t.Errorf("Expecting results from the search, got none")
	}
}
