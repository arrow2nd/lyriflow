package lrclib

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://lrclib.net/api/get"

type Response struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	TrackName    string  `json:"trackName"`
	ArtistName   string  `json:"artistName"`
	AlbumName    string  `json:"albumName"`
	Duration     float64 `json:"duration"`
	Instrumental bool    `json:"instrumental"`
	PlainLyrics  string  `json:"plainLyrics"`
	SyncedLyrics string  `json:"syncedLyrics"`
}

type Client struct {
	httpClient *http.Client
	version    string
}

func NewClient(version string) *Client {
	return &Client{
		httpClient: &http.Client{},
		version:    version,
	}
}

func (c *Client) getUserAgent() string {
	return fmt.Sprintf("LyriFlow/%s (https://github.com/arrow2nd/lyriflow)", c.version)
}

func (c *Client) GetLyrics(trackName, artistName, albumName string) (*Response, error) {
	params := url.Values{}
	params.Add("track_name", trackName)
	params.Add("artist_name", artistName)
	params.Add("album_name", albumName)

	reqURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.getUserAgent())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch lyrics: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("lyrics not found")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
