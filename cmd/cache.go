package cmd

import (
	"fmt"

	"github.com/arrow2nd/lyriflow/internal/cache"
	"github.com/urfave/cli/v2"
)

func PurgeCache(c *cli.Context) error {
	cacheStore, err := cache.NewCache()
	if err != nil {
		return fmt.Errorf("failed to initialize cache: %w", err)
	}

	if err := cacheStore.Clear(); err != nil {
		return fmt.Errorf("failed to clear cache: %w", err)
	}

	fmt.Println("Cache cleared successfully")
	return nil
}
