package cache

import (
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	aliaeState "github.com/jandedobbeleer/aliae/src/state"
)

const Version = 2

type Entry[T any] struct {
	Value        T
	SchemaHash   string
	CreatedAtUTC string
	SourceFiles  []SourceFile
	Version      int
	ComputeVars  bool
}

type SourceFile struct {
	Path         string
	Size         int64
	ModTimeNanos int64
}

func Load[T any](cachePath, schemaHash string, computeVars bool) (*T, bool, error) {
	file, err := os.Open(cachePath)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	defer file.Close()

	var entry Entry[T]
	if err := gob.NewDecoder(file).Decode(&entry); err != nil {
		return nil, false, nil
	}

	if entry.Version != Version || entry.SchemaHash != schemaHash || entry.ComputeVars != computeVars {
		return nil, false, nil
	}

	if !isSourceFingerprintValid(entry.SourceFiles) {
		return nil, false, nil
	}

	return &entry.Value, true, nil
}

func Store[T any](cachePath, schemaHash string, computeVars bool, inputs []string, value T) error {
	sourceFiles, err := fingerprintSourceFiles(inputs)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0o700); err != nil {
		return err
	}

	file, err := os.Create(cachePath)
	if err != nil {
		return err
	}
	defer file.Close()

	entry := Entry[T]{
		Version:      Version,
		SchemaHash:   schemaHash,
		ComputeVars:  computeVars,
		Value:        value,
		SourceFiles:  sourceFiles,
		CreatedAtUTC: time.Now().UTC().Format(time.RFC3339Nano),
	}

	return gob.NewEncoder(file).Encode(entry)
}

func Path(configPath, contextKey string) string {
	key := filepath.Clean(configPath) + "|" + contextKey
	sum := sha256.Sum256([]byte(key))
	fileName := fmt.Sprintf("config-cache-%s.gob", hex.EncodeToString(sum[:8]))
	return aliaeState.Path(fileName)
}

func fingerprintSourceFiles(inputs []string) ([]SourceFile, error) {
	normalized := make([]string, 0, len(inputs))
	seen := map[string]struct{}{}
	for _, input := range inputs {
		clean := filepath.Clean(strings.TrimSpace(input))
		if clean == "" {
			continue
		}
		if _, ok := seen[clean]; ok {
			continue
		}
		seen[clean] = struct{}{}
		normalized = append(normalized, clean)
	}
	slices.Sort(normalized)

	files := make([]SourceFile, 0, len(normalized))
	for _, path := range normalized {
		stat, err := os.Stat(path)
		if err != nil {
			return nil, err
		}
		files = append(files, SourceFile{
			Path:         path,
			Size:         stat.Size(),
			ModTimeNanos: stat.ModTime().UTC().UnixNano(),
		})
	}

	return files, nil
}

func isSourceFingerprintValid(files []SourceFile) bool {
	if len(files) == 0 {
		return false
	}

	for _, file := range files {
		stat, err := os.Stat(file.Path)
		if err != nil {
			return false
		}
		if stat.Size() != file.Size {
			return false
		}
		if stat.ModTime().UTC().UnixNano() != file.ModTimeNanos {
			return false
		}
	}

	return true
}
