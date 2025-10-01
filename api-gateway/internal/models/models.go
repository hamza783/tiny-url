package models

// API shorten endpoint request and responses
type ShortenURLRequest struct {
	LongUrl     string   `json:"long_url"`
	LongUrlList []string `json:"long_url_list"`
}

type ShortenURLResponse struct {
	LongUrl  string            `json:"long_url"`
	ShortUrl string            `json:"short_url"`
	BatchId  string            `json:"batch_id"`
	UrlMaps  map[string]string `json:"urls_map"`
}

type APIResponse struct {
	Data any `json:"data"`
}

// Shortening service request and response
type ShortenURLServiceBatchResponse struct {
	LongUrl  string            `json:"long_url"`
	ShortUrl string            `json:"short_url"`
	BatchId  string            `json:"batch_id"`
	UrlMaps  map[string]string `json:"urls_map"`
}

type URLShoteningServiceResponse struct {
	Data ShortenURLServiceBatchResponse `json:"data"`
}
