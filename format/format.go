// Package format handles parsing and serializing config files for supported formats.
package format

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Format represents a supported config file format.
type Format string

const (
	JSON   Format = "json"
	YAML   Format = "yaml"
	TOML   Format = "toml"
	DotEnv Format = "env"
)

// Detect infers the format from the file extension.
func Detect(filename string) (Format, error) {
	switch strings.ToLower(filepath.Ext(filename)) {
	case ".json":
		return JSON, nil
	case ".yaml", ".yml":
		return YAML, nil
	case ".toml":
		return TOML, nil
	case ".env":
		return DotEnv, nil
	default:
		return "", fmt.Errorf("unsupported extension %q (use .json, .yaml/.yml, .toml, or .env)", filepath.Ext(filename))
	}
}

// Parse decodes raw bytes into a map[string]interface{} using the given format.
func Parse(data []byte, f Format) (map[string]interface{}, error) {
	switch f {
	case JSON:
		var out map[string]interface{}
		return out, json.Unmarshal(data, &out)
	case YAML:
		var out map[string]interface{}
		return out, yaml.Unmarshal(data, &out)
	case TOML:
		var out map[string]interface{}
		_, err := toml.Decode(string(data), &out)
		return out, err
	case DotEnv:
		return parseDotenv(data)
	default:
		return nil, fmt.Errorf("unknown format: %s", f)
	}
}

// Serialize encodes a map[string]interface{} back to the given format.
func Serialize(data map[string]interface{}, f Format) ([]byte, error) {
	switch f {
	case JSON:
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return nil, err
		}
		return append(b, '\n'), nil
	case YAML:
		return yaml.Marshal(data)
	case TOML:
		var buf bytes.Buffer
		if err := toml.NewEncoder(&buf).Encode(data); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	case DotEnv:
		return serializeDotenv(data)
	default:
		return nil, fmt.Errorf("unknown format: %s", f)
	}
}
