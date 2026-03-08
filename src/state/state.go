package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jandedobbeleer/aliae/src/context"
)

type FileFormat string

const (
	FormatJSON FileFormat = "json"
	FormatText FileFormat = "text"
)

type fileState struct {
	LastRun string `json:"lastRun"`
}

type Metrics struct {
	ShouldRunCount int64
	ShouldRunTime  time.Duration
}

var (
	metricsEnabled         atomic.Bool
	shouldRunCount         atomic.Int64
	shouldRunDurationNanos atomic.Int64
)

func EnableMetrics(enabled bool) {
	metricsEnabled.Store(enabled)
}

func ResetMetrics() {
	shouldRunCount.Store(0)
	shouldRunDurationNanos.Store(0)
}

func SnapshotMetrics() Metrics {
	return Metrics{
		ShouldRunCount: shouldRunCount.Load(),
		ShouldRunTime:  time.Duration(shouldRunDurationNanos.Load()),
	}
}

func RootDir() string {
	osName := runtimeOS()
	if osName == context.WINDOWS {
		localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA"))
		if len(localAppData) == 0 {
			localAppData = filepath.Join(context.Home(), "AppData", "Local")
		}
		return filepath.Join(localAppData, "aliae", "state")
	}
	if osName == context.DARWIN {
		return filepath.Join(context.Home(), "Library", "Application Support", "aliae", "State")
	}

	return filepath.Join(context.Home(), ".local", "aliae", "state")
}

func Path(file string) string {
	return filepath.Join(RootDir(), file)
}

func IsValidFileName(file string) bool {
	name := strings.TrimSpace(file)
	if len(name) == 0 {
		return false
	}

	if filepath.Base(name) != name {
		return false
	}

	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return false
	}

	if strings.Contains(name, "..") {
		return false
	}

	return true
}

func ReadLastRun(path string) (*time.Time, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(string(data))
	if len(trimmed) == 0 {
		return nil, nil
	}

	if parsed := parseTimestamp(trimmed); parsed != nil {
		return parsed, nil
	}

	var current fileState
	if err := json.Unmarshal([]byte(trimmed), &current); err == nil {
		if parsed := parseTimestamp(current.LastRun); parsed != nil {
			return parsed, nil
		}
	}

	return nil, fmt.Errorf("invalid state timestamp format in %s", path)
}

func ShouldRun(path string, runEvery time.Duration, now time.Time) (bool, *time.Time, error) {
	if metricsEnabled.Load() {
		start := time.Now()
		run, lastRun, err := shouldRun(path, runEvery, now)
		shouldRunCount.Add(1)
		shouldRunDurationNanos.Add(time.Since(start).Nanoseconds())
		return run, lastRun, err
	}

	return shouldRun(path, runEvery, now)
}

func shouldRun(path string, runEvery time.Duration, now time.Time) (bool, *time.Time, error) {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return true, nil, nil
	}
	if err != nil {
		return true, nil, err
	}

	lastRun := info.ModTime().UTC()

	if runEvery <= 0 {
		return false, &lastRun, nil
	}

	return now.Sub(lastRun) >= runEvery, &lastRun, nil
}

func WriteLastRun(path string, now time.Time, format FileFormat) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}

	writeAndStamp := func(data []byte) error {
		if err := os.WriteFile(path, data, 0o600); err != nil {
			return err
		}

		timestamp := now.UTC()
		return os.Chtimes(path, timestamp, timestamp)
	}

	switch format {
	case "", FormatJSON:
		payload, err := json.Marshal(fileState{LastRun: now.UTC().Format(time.RFC3339Nano)})
		if err != nil {
			return err
		}
		return writeAndStamp(payload)
	case FormatText:
		return writeAndStamp([]byte(now.UTC().Format(time.RFC3339Nano)))
	default:
		return fmt.Errorf("unsupported state format: %s", format)
	}
}

func parseTimestamp(value string) *time.Time {
	formats := []string{
		time.RFC3339Nano,
		time.RFC3339,
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, strings.TrimSpace(value))
		if err == nil {
			return &parsed
		}
	}

	return nil
}

func runtimeOS() string {
	if current := context.GetCurrent(); current != nil && len(current.OS) > 0 {
		return current.OS
	}

	return runtime.GOOS
}
