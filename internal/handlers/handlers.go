package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"youtube-api/internal/cache"
	"youtube-api/internal/config"
	"youtube-api/internal/youtube"
)

type Handler struct {
	config   *config.Config
	cache    *cache.Cache
	ytClient *youtube.Service
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(cfg *config.Config, cache *cache.Cache, ytClient *youtube.Service) *Handler {
	return &Handler{
		config:   cfg,
		cache:    cache,
		ytClient: ytClient,
	}
}

func (h *Handler) GetVideos(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel_id")
	if channelID == "" {
		h.writeError(w, http.StatusBadRequest, "channel_id is required")
		return
	}

	if !h.config.IsChannelAllowed(channelID) {
		h.writeError(w, http.StatusForbidden, "channel not allowed")
		return
	}

	maxResults := h.parseMaxResults(r.URL.Query().Get("max_results"))

	cacheKey := cache.VideosKey(channelID)
	var videos []youtube.Video

	found, err := h.cache.Get(r.Context(), cacheKey, &videos)
	if err != nil {
		log.Printf("cache get error: %v", err)
	}

	if found {
		h.writeJSON(w, http.StatusOK, videos)
		return
	}

	videos, err = h.ytClient.GetLatestVideos(r.Context(), channelID, maxResults)
	if err != nil {
		log.Printf("youtube api error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to fetch videos")
		return
	}

	if err := h.cache.Set(r.Context(), cacheKey, videos); err != nil {
		log.Printf("cache set error: %v", err)
	}

	h.writeJSON(w, http.StatusOK, videos)
}

func (h *Handler) GetPlaylists(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query().Get("channel_id")
	if channelID == "" {
		h.writeError(w, http.StatusBadRequest, "channel_id is required")
		return
	}

	if !h.config.IsChannelAllowed(channelID) {
		h.writeError(w, http.StatusForbidden, "channel not allowed")
		return
	}

	maxResults := h.parseMaxResults(r.URL.Query().Get("max_results"))

	cacheKey := cache.PlaylistsKey(channelID)
	var playlists []youtube.Playlist

	found, err := h.cache.Get(r.Context(), cacheKey, &playlists)
	if err != nil {
		log.Printf("cache get error: %v", err)
	}

	if found {
		h.writeJSON(w, http.StatusOK, playlists)
		return
	}

	playlists, err = h.ytClient.GetLatestPlaylists(r.Context(), channelID, maxResults)
	if err != nil {
		log.Printf("youtube api error: %v", err)
		h.writeError(w, http.StatusInternalServerError, "failed to fetch playlists")
		return
	}

	if err := h.cache.Set(r.Context(), cacheKey, playlists); err != nil {
		log.Printf("cache set error: %v", err)
	}

	h.writeJSON(w, http.StatusOK, playlists)
}

func (h *Handler) parseMaxResults(value string) int {
	if value == "" {
		return 10
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return 10
	}
	if n > 50 {
		return 50
	}
	return n
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("write response error: %v", err)
	}
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{Error: message})
}
