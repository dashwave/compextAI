package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/burnerlee/compextAI/internal/logger"
	"gorm.io/gorm"
)

type executorClient struct {
	BaseURL string
}

func getExecutorClient() *executorClient {
	return &executorClient{
		BaseURL: os.Getenv("EXECUTOR_BASE_URL"),
	}
}

type ExecuteParams struct {
	Timeout time.Duration
}

func (c *executorClient) getRequest(execRoute, method string, data interface{}) (*http.Request, error) {
	var body io.Reader
	if data != nil {
		dataJson, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("error marshalling data: %w", err)
		}
		body = bytes.NewBuffer(dataJson)
	}
	request, err := http.NewRequest(method, c.BaseURL+execRoute, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	return request, nil
}

func Execute(db *gorm.DB, execRoute string, executeParams *ExecuteParams, threadExecutionData interface{}, threadExecutionIdentifier string, messages interface{}) (int, interface{}, error) {
	executorClient := getExecutorClient()

	// update thread execution metadata
	if err := UpdateThreadExecutionMetadata(db, threadExecutionIdentifier, threadExecutionData, messages); err != nil {
		logger.GetLogger().Errorf("Error updating thread execution metadata: %v", err)
		return -1, nil, err
	}

	request, err := executorClient.getRequest(execRoute, "POST", threadExecutionData)
	if err != nil {
		return -1, nil, fmt.Errorf("error getting request: %w", err)
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: executeParams.Timeout,
	}

	response, err := client.Do(request)
	if err != nil {
		return -1, nil, fmt.Errorf("error executing request: %w", err)
	}

	defer response.Body.Close()

	var responseData interface{}
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return response.StatusCode, nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.StatusCode, responseData, nil
}
