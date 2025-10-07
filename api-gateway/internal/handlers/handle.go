package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	models "github.com/hamza4253/tiny-url/gateway/internal/models"
	"github.com/hamza4253/tiny-url/gateway/internal/publisher"
	services "github.com/hamza4253/tiny-url/gateway/internal/services"
	pb "github.com/hamza4253/tiny-url/shared/api/gen"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type Handler struct {
	shortenerClient   *services.URLShortnerClient
	redirectionClient pb.RedirectionServiceClient
	publisher         *publisher.ShorteningPublisher
}

func NewHandler(s *services.URLShortnerClient, c pb.RedirectionServiceClient, p *publisher.ShorteningPublisher) *Handler {
	return &Handler{
		shortenerClient:   s,
		redirectionClient: c,
		publisher:         p,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/shorten", h.ShortenURL)
	mux.HandleFunc("POST /api/shorten/all", h.ShortenURLs)
	mux.HandleFunc("GET /api/{short_url}", h.LookupUrl)
	mux.HandleFunc("GET /api/all/{batch_id}", h.FetchUrlsByBatchId)
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody models.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	longUrl := reqBody.LongUrl
	if longUrl == "" {
		http.Error(w, "Invalid request body. long_url is required", http.StatusBadRequest)
		return
	}

	fmt.Println("Calling shorten service client ShortenURL using REST")
	short_url, err := h.shortenerClient.ShortenURL(ctx, longUrl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	urlsResp := models.ShortenURLResponse{
		ShortUrl: short_url,
		LongUrl:  longUrl,
	}
	response := models.APIResponse{
		Data: urlsResp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) ShortenURLs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	batchId, err := createShortRandomUrl()
	if err != nil {
		http.Error(w, "Unexpected error occurred", http.StatusInternalServerError)
		return
	}

	var reqBody models.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	longUrlList := reqBody.LongUrlList

	for _, longUrl := range longUrlList {
		url := strings.TrimSpace(longUrl)
		h.publisher.Publish(ctx, batchId, url)
	}

	urlsResp := models.ShortenURLResponse{
		BatchId: batchId,
	}
	response := models.APIResponse{
		Data: urlsResp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// V2 iteration. Calls url-redirection-service microservice using gRPC to get full URL
func (h *Handler) LookupUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get short url from url
	shortUrl := r.PathValue("short_url")

	request := pb.LookupRequest{
		ShortUrl: shortUrl,
	}
	fmt.Println("Calling redirection service client LookupURL using gRPC")
	longURLResponse, err := h.redirectionClient.LookupURL(ctx, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	longURL := longURLResponse.LongUrl
	if longURL == "" {
		http.Error(w, "no long URL found for given short URL", http.StatusNotFound)
		return
	}

	// Ensure the URL has a protocol prefix
	if !strings.HasPrefix(longURL, "http://") && !strings.HasPrefix(longURL, "https://") {
		longURL = "https://" + longURL
	}
	http.Redirect(w, r, longURL, http.StatusFound)
}

func (h *Handler) FetchUrlsByBatchId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get short url from url
	batchId := r.PathValue("batch_id")

	urlsMap, err := h.shortenerClient.FetchURLsByBatchId(ctx, batchId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	urlsResp := models.ShortenURLResponse{
		UrlMaps: urlsMap,
	}
	response := models.APIResponse{
		Data: urlsResp,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createShortRandomUrl() (string, error) {
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	id, err := gonanoid.Generate(alphabet, 6)
	if err != nil {
		return "", err
	}

	return id, nil
}
