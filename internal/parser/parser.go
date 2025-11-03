package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type LyricLine struct {
	Time float64
	Text string
}

var lrcRegex = regexp.MustCompile(`\[(\d+):(\d+)\.(\d+)\]\s*(.*)`)

func ParseLRC(syncedLyrics string) ([]LyricLine, error) {
	lines := strings.Split(syncedLyrics, "\n")
	var result []LyricLine

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		matches := lrcRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		minutes, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		seconds, err := strconv.Atoi(matches[2])
		if err != nil {
			continue
		}

		centiseconds, err := strconv.Atoi(matches[3])
		if err != nil {
			continue
		}

		timeInSeconds := float64(minutes)*60 + float64(seconds) + float64(centiseconds)/100.0
		text := strings.TrimSpace(matches[4])

		result = append(result, LyricLine{
			Time: timeInSeconds,
			Text: text,
		})
	}

	if len(result) == 0 {
		return nil, errors.New("no valid lyrics found")
	}

	return result, nil
}

func GetLyricAtTime(lines []LyricLine, duration float64) (string, error) {
	if len(lines) == 0 {
		return "", errors.New("no lyrics available")
	}

	var currentLine *LyricLine

	for i := range lines {
		if lines[i].Time > duration {
			break
		}
		currentLine = &lines[i]
	}

	if currentLine == nil {
		return "", fmt.Errorf("no lyrics at time %.2f", duration)
	}

	return currentLine.Text, nil
}
