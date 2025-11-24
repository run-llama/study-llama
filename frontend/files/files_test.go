package files

import (
	"os"
	"testing"
)

func TestUploadFile(t *testing.T) {
	if _, ok := os.LookupEnv("LLAMA_CLOUD_API_KEY"); !ok {
		t.Skip("LLAMA_CLOUD_API_KEY not available")
	}
	file := "../testfiles/the-future-of-vibe-coding.pdf"
	src, _ := os.Open(file)
	_, err := UploadFile(src, file)
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error())
	}
}
