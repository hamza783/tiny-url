package handler

import (
	"encoding/json"
	"net/http"

	"github.com/hamza4253/tiny-url/shortener/internal/helpers"
	"github.com/hamza4253/tiny-url/shortener/internal/models"

	shorten "github.com/hamza4253/tiny-url/shortener/internal/service"
)

type Handler struct {
	service shorten.ShortenService
}

func NewHandler(shortenService *shorten.ShortenService) *Handler {
	return &Handler{
		service: *shortenService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/urls/shorten", h.ShortenURL)
	mux.HandleFunc("GET /api/urls/{short_url}", h.GetLongUrl)
	mux.HandleFunc("GET /api/urls/all/{batch_id}", h.GetUrlsByBatchId)
}

// Calls ShortenURL service method and creates response for ShortenURL endpoint
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var reqBody models.RequestBody
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// shorten url
	shortUrl, err := h.service.Shorten(ctx, reqBody.LongUrl)
	if err != nil {
		http.Error(w, "Error shortening the url", http.StatusInternalServerError)
		return
	}

	// create a response object
	urlResponse := models.URLResponse{
		ShortUrl: shortUrl,
		LongUrl:  reqBody.LongUrl,
	}
	response := models.APIResponse{
		Data: urlResponse,
	}
	helpers.WriteResponse(w, response, http.StatusOK)
}

// Calls GetLongUrl service method and creates response for GetLongUrl endpoint.
// Used only for testing. Real application uses url-redirection-service
func (h *Handler) GetLongUrl(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	shortUrl := r.PathValue("short_url")
	// get long url
	longUrl, err := h.service.GetFullURL(ctx, shortUrl)
	if err != nil {
		http.Error(w, "Error getting the long url", http.StatusInternalServerError)
		return
	}

	// create a response object
	response := models.APIResponse{
		Data: longUrl,
	}

	helpers.WriteResponse(w, response, http.StatusOK)
}

func (h *Handler) GetUrlsByBatchId(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	batchId := r.PathValue("batch_id")
	// get long url
	urlsMap, err := h.service.GetUrlsByBatchId(ctx, batchId)
	if err != nil {
		http.Error(w, "Error getting the long url", http.StatusInternalServerError)
		return
	}

	// create a response object
	urlResponse := models.URLResponse{
		UrlsMap: urlsMap,
	}
	response := models.APIResponse{
		Data: urlResponse,
	}
	helpers.WriteResponse(w, response, http.StatusOK)
}
