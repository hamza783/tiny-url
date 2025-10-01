package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/hamza4253/tiny-url/gateway/internal/models"
)

type URLShortnerClient struct {
	BaseURL    string
	httpClient *http.Client
}

func NewURLShorteningClient(baseUrl string) *URLShortnerClient {
	return &URLShortnerClient{
		BaseURL: baseUrl,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *URLShortnerClient) ShortenURL(ctx context.Context, longUrl string) (string, error) {
	url := fmt.Sprintf("%s/api/urls/shorten", c.BaseURL)
	// Create request
	payload := map[string]string{"long_url": longUrl}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("failed to create request payload")
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonPayload))
	if err != nil {
		return "", err
	}
	// Add Header
	request.Header.Add("Content-Type", "application/json")

	// Get Response
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to shorten url")
	}

	var response models.URLShoteningServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	data := response.Data
	return data.ShortUrl, nil
}

// TODO: replace with url-redirection-service client
func (c *URLShortnerClient) LookupURL(ctx context.Context, shortUrl string) (string, error) {
	url := fmt.Sprintf("%s/api/urls/%s", c.BaseURL, shortUrl)
	fmt.Println("url shortening service LookupURL")
	// Create request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	// Get Response
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error ==>", resp.StatusCode)
		return "", errors.New("failed to get long url")
	}
	fmt.Println("url shortening service LookupURL status 2", resp.StatusCode)
	var response map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}
	fmt.Println("Long url ====>", response["data"])
	return response["data"], nil
}

func (c *URLShortnerClient) FetchURLsByBatchId(ctx context.Context, batchId string) (map[string]string, error) {
	url := fmt.Sprintf("%s/api/urls/all/%s", c.BaseURL, batchId)
	// Create request
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Add Header
	request.Header.Add("Content-Type", "application/json")

	// Get Response
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get urls by batch id")
	}

	var response models.URLShoteningServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	data := response.Data
	return data.UrlMaps, nil
}
