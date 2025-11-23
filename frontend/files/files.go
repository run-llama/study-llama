package files

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

type UploadedFile struct {
	CreatedAt      *string                `json:"created_at"`
	DataSourceID   *string                `json:"data_source_id"`
	ExternalFileID *string                `json:"external_file_id"`
	FileSize       *int64                 `json:"file_size"`
	FileType       *string                `json:"file_type"`
	ID             string                 `json:"id"`
	LastModifiedAt *string                `json:"last_modified_at"`
	Name           string                 `json:"name"`
	PermissionInfo map[string]interface{} `json:"permission_info"`
	ProjectID      string                 `json:"project_id"`
	ResourceInfo   map[string]interface{} `json:"resource_info"`
	UpdatedAt      *string                `json:"updated_at"`
}

func UploadFile(file io.Reader, fileName string) (string, error) {
	apiKey := os.Getenv("LLAMA_CLOUD_API_KEY")
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	fileWriter, _ := writer.CreateFormFile("upload_file", fileName)

	io.Copy(fileWriter, file)

	writer.Close()
	url := "https://api.cloud.llamaindex.ai/api/v1/files"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, &requestBody)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "multipart/form-data")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiKey)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var fl UploadedFile
	err = json.Unmarshal(body, &fl)
	if err != nil {
		return "", err
	}
	return fl.ID, nil
}
