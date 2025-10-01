package models

// request/response structs
type RequestBody struct {
	LongUrl string `json:"long_url"`
}

type APIResponse struct {
	Data any `json:"data"`
}

type URLResponse struct {
	ShortUrl string            `json:"short_url"`
	LongUrl  string            `json:"long_url"`
	UrlsMap  map[string]string `json:"urls_map"`
}
