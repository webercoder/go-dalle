package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	DefaultURL = "https://api.openai.com/v1/images/generations"
)

type Model string

var (
	ModelDallE3 Model = "dall-e-3"
)

type Quality string

var (
	QualityHD       Quality = "hd"
	QualityStandard Quality = "standard"
)

type Size string

var (
	Size1024x1024 Size = "1024x1024"
	Size1024x1792 Size = "1024x1792"
	Size1792x1024 Size = "1792x1024"
)

type DallERequest struct {
	Model   Model   `json:"model"`
	Prompt  string  `json:"prompt"`
	Size    Size    `json:"size"`
	Quality Quality `json:"quality"`
	Count   int     `json:"n"`
}

type DallEResponse struct {
	Created int64               `json:"created"`
	Data    []DallEResponseData `json:"data"`
}

type DallEResponseData struct {
	RevisedPrompt string `json:"revised_prompt"`
	URL           string `json:"url"`
}

type DallEClient struct {
	apiKey string
	client *http.Client
	url    string
}

func NewDallEClient(key string) *DallEClient {
	return &DallEClient{
		apiKey: key,
		client: &http.Client{},
		url:    DefaultURL,
	}
}

func (c *DallEClient) Request(ctx context.Context, req DallERequest) (*DallEResponse, error) {
	if req.Count == 0 {
		req.Count = 1
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(req); err != nil {
		return nil, err
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, &buf)
	if err != nil {
		return nil, err
	}

	hreq.Header.Add("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))
	hreq.Header.Add("Content-Type", "application/json")

	resp, err := c.client.Do(hreq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		if bytes, err := io.ReadAll(resp.Body); err == nil {
			return nil, fmt.Errorf("client returned status %v: %v", resp.StatusCode, string(bytes))
		}

		return nil, fmt.Errorf("client returned status %v", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dresp DallEResponse
	if err := json.Unmarshal(bytes, &dresp); err != nil {
		return nil, err
	}

	return &dresp, nil
}
