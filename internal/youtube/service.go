package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	baseURL           = "https://www.googleapis.com/youtube/v3"
	defaultMaxResults = 10
)

type Service struct {
	apiKey     string
	httpClient *http.Client
}

type Video struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PublishedAt  time.Time `json:"published_at"`
	ThumbnailURL string    `json:"thumbnail_url"`
	ChannelID    string    `json:"channel_id"`
	ChannelTitle string    `json:"channel_title"`
}

type Playlist struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	PublishedAt  time.Time `json:"published_at"`
	ThumbnailURL string    `json:"thumbnail_url"`
	ChannelID    string    `json:"channel_id"`
	ChannelTitle string    `json:"channel_title"`
	ItemCount    int       `json:"item_count"`
}

type searchResponse struct {
	Items []struct {
		ID struct {
			Kind       string `json:"kind"`
			VideoID    string `json:"videoId"`
			PlaylistID string `json:"playlistId"`
		} `json:"id"`
		Snippet struct {
			Title        string    `json:"title"`
			Description  string    `json:"description"`
			PublishedAt  time.Time `json:"publishedAt"`
			ChannelID    string    `json:"channelId"`
			ChannelTitle string    `json:"channelTitle"`
			Thumbnails   struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
	} `json:"items"`
}

type playlistsResponse struct {
	Items []struct {
		ID      string `json:"id"`
		Snippet struct {
			Title        string    `json:"title"`
			Description  string    `json:"description"`
			PublishedAt  time.Time `json:"publishedAt"`
			ChannelID    string    `json:"channelId"`
			ChannelTitle string    `json:"channelTitle"`
			Thumbnails   struct {
				Default struct {
					URL string `json:"url"`
				} `json:"default"`
				Medium struct {
					URL string `json:"url"`
				} `json:"medium"`
				High struct {
					URL string `json:"url"`
				} `json:"high"`
			} `json:"thumbnails"`
		} `json:"snippet"`
		ContentDetails struct {
			ItemCount int `json:"itemCount"`
		} `json:"contentDetails"`
	} `json:"items"`
}

func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Service) GetLatestVideos(ctx context.Context, channelID string, maxResults int) ([]Video, error) {
	if maxResults <= 0 {
		maxResults = defaultMaxResults
	}

	params := url.Values{}
	params.Set("part", "snippet")
	params.Set("channelId", channelID)
	params.Set("type", "video")
	params.Set("order", "date")
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	params.Set("key", s.apiKey)

	reqURL := fmt.Sprintf("%s/search?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube api returned status %d", resp.StatusCode)
	}

	var searchResp searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	videos := make([]Video, 0, len(searchResp.Items))
	for _, item := range searchResp.Items {
		thumbnail := item.Snippet.Thumbnails.High.URL
		if thumbnail == "" {
			thumbnail = item.Snippet.Thumbnails.Medium.URL
		}
		if thumbnail == "" {
			thumbnail = item.Snippet.Thumbnails.Default.URL
		}

		videos = append(videos, Video{
			ID:           item.ID.VideoID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  item.Snippet.PublishedAt,
			ThumbnailURL: thumbnail,
			ChannelID:    item.Snippet.ChannelID,
			ChannelTitle: item.Snippet.ChannelTitle,
		})
	}

	return videos, nil
}

func (s *Service) GetLatestPlaylists(ctx context.Context, channelID string, maxResults int) ([]Playlist, error) {
	if maxResults <= 0 {
		maxResults = defaultMaxResults
	}

	params := url.Values{}
	params.Set("part", "snippet,contentDetails")
	params.Set("channelId", channelID)
	params.Set("maxResults", fmt.Sprintf("%d", maxResults))
	params.Set("key", s.apiKey)

	reqURL := fmt.Sprintf("%s/playlists?%s", baseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("youtube api returned status %d", resp.StatusCode)
	}

	var playlistsResp playlistsResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlistsResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	playlists := make([]Playlist, 0, len(playlistsResp.Items))
	for _, item := range playlistsResp.Items {
		thumbnail := item.Snippet.Thumbnails.High.URL
		if thumbnail == "" {
			thumbnail = item.Snippet.Thumbnails.Medium.URL
		}
		if thumbnail == "" {
			thumbnail = item.Snippet.Thumbnails.Default.URL
		}

		playlists = append(playlists, Playlist{
			ID:           item.ID,
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  item.Snippet.PublishedAt,
			ThumbnailURL: thumbnail,
			ChannelID:    item.Snippet.ChannelID,
			ChannelTitle: item.Snippet.ChannelTitle,
			ItemCount:    item.ContentDetails.ItemCount,
		})
	}

	return playlists, nil
}
