package cmd

import (
	"fmt"

	"github.com/arrow2nd/lyriflow/internal/cache"
	"github.com/arrow2nd/lyriflow/internal/lrclib"
	"github.com/arrow2nd/lyriflow/internal/parser"
	"github.com/urfave/cli/v2"
)

func GetLyrics(version string) cli.ActionFunc {
	return func(c *cli.Context) error {
		title := c.String("title")
		artist := c.String("artist")
		album := c.String("album")
		position := c.Float64("position")

		cacheStore, err := cache.NewCache()
		if err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}

		// ファイルロック取得
		fileLock, err := cacheStore.Lock(title, artist, album)
		if err != nil {
			return fmt.Errorf("failed to acquire lock: %w", err)
		}
		defer fileLock.Unlock()

		var lines []parser.LyricLine

		cached, err := cacheStore.Get(title, artist, album)
		if err != nil {
			return fmt.Errorf("failed to get cache: %w", err)
		}

		if cached != nil {
			// キャッシュヒット
			if cached.NotFound {
				fmt.Println("Lyrics not found")
				return nil
			}
			lines = cached.ParsedLyrics
		} else {
			// キャッシュミス - API取得
			client := lrclib.NewClient(version)
			response, err := client.GetLyrics(title, artist, album)
			if err != nil {
				if err.Error() == "lyrics not found" {
					// 失敗をキャッシュ
					if err := cacheStore.SetNotFound(title, artist, album); err != nil {
						return fmt.Errorf("failed to save cache: %w", err)
					}
					fmt.Println("Lyrics not found")
					return nil
				}
				return fmt.Errorf("failed to fetch lyrics: %w", err)
			}

			if response.Instrumental {
				// インスト曲も失敗としてキャッシュ
				if err := cacheStore.SetNotFound(title, artist, album); err != nil {
					return fmt.Errorf("failed to save cache: %w", err)
				}
				fmt.Println("No lyrics available")
				return nil
			}

			if response.SyncedLyrics == "" {
				// 歌詞なしも失敗としてキャッシュ
				if err := cacheStore.SetNotFound(title, artist, album); err != nil {
					return fmt.Errorf("failed to save cache: %w", err)
				}
				fmt.Println("No lyrics available")
				return nil
			}

			// パース
			lines, err = parser.ParseLRC(response.SyncedLyrics)
			if err != nil {
				return fmt.Errorf("failed to parse lyrics: %w", err)
			}

			// パース済みデータをキャッシュ
			if err := cacheStore.Set(title, artist, album, response, lines); err != nil {
				return fmt.Errorf("failed to save cache: %w", err)
			}
		}

		lyric, err := parser.GetLyricAtTime(lines, position)
		if err != nil {
			fmt.Println("No lyrics available")
			return nil
		}

		if lyric == "" {
			fmt.Println("(instrumental)")
			return nil
		}

		fmt.Println(lyric)
		return nil
	}
}
