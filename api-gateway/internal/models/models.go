package models

type ShortenURLRequest struct {
	LongUrl string `json:"long_url"`
}

type ShortenURLResponse struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

type APIResponse struct {
	Data any `json:"data"`
}

type URLShoteningServiceResponse struct {
	Data ShortenURLResponse `json:"data"`
}
