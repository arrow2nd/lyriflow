package main

import (
	"log"
	"os"

	"github.com/arrow2nd/lyriflow/cmd"
	"github.com/urfave/cli/v2"
)

const Version = "0.0.2"

func main() {
	app := &cli.App{
		Name:  "lyriflow",
		Usage: "Get synchronized lyrics for your music",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get lyrics at specified time",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "title",
						Aliases:  []string{"t"},
						Usage:    "Track title",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "artist",
						Aliases:  []string{"a"},
						Usage:    "Artist name",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "album",
						Aliases:  []string{"A"},
						Usage:    "Album name",
						Required: true,
					},
					&cli.Float64Flag{
						Name:     "position",
						Aliases:  []string{"p"},
						Usage:    "Current playback position in seconds",
						Required: true,
					},
					&cli.BoolFlag{
						Name:  "waybar",
						Usage: "Output in waybar JSON format",
					},
				},
				Action: cmd.GetLyrics(Version),
			},
			{
				Name:   "cache-purge",
				Usage:  "Clear lyrics cache",
				Action: cmd.PurgeCache,
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Show version information",
				Action:  cmd.ShowVersion(Version),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
