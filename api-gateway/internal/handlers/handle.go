package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	models "github.com/hamza4253/tiny-url/gateway/internal/models"
	services "github.com/hamza4253/tiny-url/gateway/internal/services"
	pb "github.com/hamza4253/tiny-url/shared/api/gen"
)

type Handler struct {
	shortenerClient *services.URLShortnerClient
	client          pb.RedirectionServiceClient
}

func NewHandler(s *services.URLShortnerClient, c pb.RedirectionServiceClient) *Handler {
	return &Handler{
		shortenerClient: s,
		client:          c,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/shorten", h.ShortenURL)
	mux.HandleFunc("GET /api/{short_url}", h.LookupUrl)
}

func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody models.ShortenURLRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	longUrl := reqBody.LongUrl

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

// V2 iteration. Calls url-redirection-service microservice using gRPC to get full URL
func (h *Handler) LookupUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// get short url from url
	shortUrl := r.PathValue("short_url")

	request := pb.LookupRequest{
		ShortUrl: shortUrl,
	}
	fmt.Println("Calling redirection service client LookupURL using gRPC")
	longURLResponse, err := h.client.LookupURL(ctx, &request)
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
