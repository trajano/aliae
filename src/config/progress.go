package config

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
)

type Progress struct {
	Enabled         bool                  `yaml:"-"`
	StartPercentage float64               `yaml:"start_percentage,omitempty"`
	Internal        float64               `yaml:"internal,omitempty"`
	EndPercentage   ProgressEndPercentage `yaml:"end_percentage,omitempty"`
}

type ProgressEndPercentage struct {
	Reset bool
	Value float64
}

func (p *Progress) UnmarshalYAML(data []byte) error {
	return progressUnmarshaler(p, data)
}

func (p Progress) MarshalYAML() (any, error) {
	if !p.Enabled {
		return false, nil
	}

	return struct {
		StartPercentage float64               `yaml:"start_percentage,omitempty"`
		Internal        float64               `yaml:"internal,omitempty"`
		EndPercentage   ProgressEndPercentage `yaml:"end_percentage,omitempty"`
	}{
		StartPercentage: p.StartPercentage,
		Internal:        p.Internal,
		EndPercentage:   p.EndPercentage,
	}, nil
}

func (p ProgressEndPercentage) MarshalYAML() (any, error) {
	if p.Reset {
		return "reset", nil
	}

	return p.Value, nil
}

func progressUnmarshaler(progress *Progress, data []byte) error {
	var raw any
	if err := yaml.NewDecoder(bytes.NewBuffer(data)).Decode(&raw); err != nil {
		return err
	}

	parsed, err := parseProgressAny(raw)
	if err != nil {
		return err
	}
	*progress = parsed
	return nil
}

func hydrateProgressFromYAML(aliae *Aliae, data []byte) error {
	if aliae == nil {
		return nil
	}

	document := map[string]any{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return err
	}

	progressValue, hasProgress := document["progress"]
	if !hasProgress {
		aliae.Progress = Progress{}
		return nil
	}

	parsed, err := parseProgressAny(progressValue)
	if err != nil {
		return err
	}

	aliae.Progress = parsed
	return nil
}

func parseProgressAny(value any) (Progress, error) {
	if value == nil {
		return Progress{}, nil
	}

	switch v := value.(type) {
	case bool:
		if v {
			return Progress{
				Enabled:         true,
				StartPercentage: 0,
				Internal:        0,
				EndPercentage: ProgressEndPercentage{
					Reset: true,
					Value: 100,
				},
			}, nil
		}
		return Progress{}, nil
	case map[string]any:
		return parseProgressObject(v)
	case map[any]any:
		typed := make(map[string]any, len(v))
		for key, rawValue := range v {
			textKey, ok := key.(string)
			if !ok {
				return Progress{}, fmt.Errorf("invalid progress key type %T", key)
			}
			typed[textKey] = rawValue
		}
		return parseProgressObject(typed)
	default:
		text := strings.TrimSpace(fmt.Sprintf("%v", v))
		if strings.EqualFold(text, "false") {
			return Progress{}, nil
		}
		if strings.EqualFold(text, "true") {
			return Progress{
				Enabled:         true,
				StartPercentage: 0,
				Internal:        0,
				EndPercentage: ProgressEndPercentage{
					Reset: true,
					Value: 100,
				},
			}, nil
		}
		return Progress{}, fmt.Errorf("progress must be false or an object")
	}
}

func parseProgressObject(raw map[string]any) (Progress, error) {
	parsed := Progress{
		Enabled:         true,
		StartPercentage: 0,
		Internal:        0,
		EndPercentage: ProgressEndPercentage{
			Reset: true,
			Value: 100,
		},
	}

	if rawStart, ok := raw["start_percentage"]; ok {
		start, err := parseProgressNumber(rawStart)
		if err != nil {
			return Progress{}, fmt.Errorf("invalid progress.start_percentage: %w", err)
		}
		parsed.StartPercentage = start
	}

	if rawInternal, ok := raw["internal"]; ok {
		internal, err := parseProgressNumber(rawInternal)
		if err != nil {
			return Progress{}, fmt.Errorf("invalid progress.internal: %w", err)
		}
		parsed.Internal = internal
	}

	if rawEnd, ok := raw["end_percentage"]; ok {
		endValue, reset, err := parseProgressEndPercentage(rawEnd)
		if err != nil {
			return Progress{}, err
		}

		parsed.EndPercentage = ProgressEndPercentage{
			Reset: reset,
			Value: endValue,
		}
	}

	return parsed, nil
}

func parseProgressNumber(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int8:
		return float64(v), nil
	case int16:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint8:
		return float64(v), nil
	case uint16:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
		if err != nil {
			return 0, fmt.Errorf("invalid number %q", v)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("invalid type %T", value)
	}
}

func parseProgressEndPercentage(value any) (float64, bool, error) {
	if v, ok := value.(string); ok {
		text := strings.TrimSpace(v)
		if strings.EqualFold(text, "reset") {
			return 100, true, nil
		}

		parsed, err := parseProgressNumber(text)
		if err != nil {
			return 0, false, fmt.Errorf("invalid progress.end_percentage value %q", v)
		}
		return parsed, false, nil
	}

	parsed, err := parseProgressNumber(value)
	if err != nil {
		return 0, false, fmt.Errorf("invalid progress.end_percentage type %T", value)
	}

	return parsed, false, nil
}
