package cache

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"strconv"

	"github.com/arrow2nd/lyriflow/internal/lrclib"
	"github.com/arrow2nd/lyriflow/internal/parser"
	"github.com/gofrs/flock"
)

const dirName = "lyriflow"

type CachedData struct {
	Response     *lrclib.Response     `json:"response"`
	ParsedLyrics []parser.LyricLine   `json:"parsed_lyrics"`
	NotFound     bool                 `json:"not_found"` // 検索失敗フラグ
}

type Cache struct {
	dir string
}

func NewCache() (*Cache, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	dir := filepath.Join(cacheDir, dirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{dir: dir}, nil
}

func (c *Cache) generateKey(trackName, artistName, albumName string) string {
	data := fmt.Sprintf("%s_%s_%s", trackName, artistName, albumName)
	h := fnv.New64a()
	h.Write([]byte(data))
	return strconv.FormatUint(h.Sum64(), 16)
}

func (c *Cache) getCachePath(trackName, artistName, albumName string) string {
	key := c.generateKey(trackName, artistName, albumName)
	return filepath.Join(c.dir, key+".json")
}

func (c *Cache) getLockPath(trackName, artistName, albumName string) string {
	key := c.generateKey(trackName, artistName, albumName)
	return filepath.Join(c.dir, key+".lock")
}

func (c *Cache) Lock(trackName, artistName, albumName string) (*flock.Flock, error) {
	lockPath := c.getLockPath(trackName, artistName, albumName)
	fileLock := flock.New(lockPath)

	if err := fileLock.Lock(); err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return fileLock, nil
}

func (c *Cache) Get(trackName, artistName, albumName string) (*CachedData, error) {
	path := c.getCachePath(trackName, artistName, albumName)

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to read cache: %w", err)
	}

	var cached CachedData
	if err := json.Unmarshal(data, &cached); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cache: %w", err)
	}

	return &cached, nil
}

func (c *Cache) Set(trackName, artistName, albumName string, response *lrclib.Response, parsedLyrics []parser.LyricLine) error {
	path := c.getCachePath(trackName, artistName, albumName)

	cached := CachedData{
		Response:     response,
		ParsedLyrics: parsedLyrics,
		NotFound:     false,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

func (c *Cache) SetNotFound(trackName, artistName, albumName string) error {
	path := c.getCachePath(trackName, artistName, albumName)

	cached := CachedData{
		Response:     nil,
		ParsedLyrics: nil,
		NotFound:     true,
	}

	data, err := json.Marshal(cached)
	if err != nil {
		return fmt.Errorf("failed to marshal cache: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache: %w", err)
	}

	return nil
}

func (c *Cache) Clear() error {
	if err := os.RemoveAll(c.dir); err != nil {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}

	if err := os.MkdirAll(c.dir, 0755); err != nil {
		return fmt.Errorf("failed to recreate cache directory: %w", err)
	}

	return nil
}
