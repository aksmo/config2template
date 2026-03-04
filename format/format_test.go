package format

import (
	"strings"
	"testing"
)

var sampleMap = map[string]interface{}{
	"host":    "${HOST}",
	"port":    float64(5432),
	"enabled": true,
}

func TestDetect(t *testing.T) {
	cases := []struct {
		file    string
		want    Format
		wantErr bool
	}{
		{"config.json", JSON, false},
		{"config.yaml", YAML, false},
		{"config.yml", YAML, false},
		{"config.toml", TOML, false},
		{"config.ini", "", true},
		{"config", "", true},
	}
	for _, c := range cases {
		got, err := Detect(c.file)
		if c.wantErr {
			if err == nil {
				t.Errorf("Detect(%q): expected error, got nil", c.file)
			}
			continue
		}
		if err != nil {
			t.Errorf("Detect(%q): unexpected error: %v", c.file, err)
		}
		if got != c.want {
			t.Errorf("Detect(%q) = %q, want %q", c.file, got, c.want)
		}
	}
}

func TestParseAndSerialize_JSON(t *testing.T) {
	input := []byte(`{"api_key": "secret", "port": 8080, "debug": false}`)
	m, err := Parse(input, JSON)
	if err != nil {
		t.Fatalf("Parse JSON: %v", err)
	}
	if m["api_key"] != "secret" {
		t.Errorf("got %v", m["api_key"])
	}
	if m["port"] != float64(8080) {
		t.Errorf("got %v", m["port"])
	}

	out, err := Serialize(sampleMap, JSON)
	if err != nil {
		t.Fatalf("Serialize JSON: %v", err)
	}
	if !strings.Contains(string(out), `"${HOST}"`) {
		t.Errorf("expected placeholder in JSON output, got:\n%s", out)
	}
}

func TestParseAndSerialize_YAML(t *testing.T) {
	input := []byte("api_key: secret\nport: 8080\ndebug: false\n")
	m, err := Parse(input, YAML)
	if err != nil {
		t.Fatalf("Parse YAML: %v", err)
	}
	if m["api_key"] != "secret" {
		t.Errorf("got %v", m["api_key"])
	}
	if m["port"] != 8080 {
		t.Errorf("got %v", m["port"])
	}

	out, err := Serialize(sampleMap, YAML)
	if err != nil {
		t.Fatalf("Serialize YAML: %v", err)
	}
	if !strings.Contains(string(out), "${HOST}") {
		t.Errorf("expected placeholder in YAML output, got:\n%s", out)
	}
}

func TestParseAndSerialize_TOML(t *testing.T) {
	input := []byte("api_key = \"secret\"\nport = 8080\ndebug = false\n")
	m, err := Parse(input, TOML)
	if err != nil {
		t.Fatalf("Parse TOML: %v", err)
	}
	if m["api_key"] != "secret" {
		t.Errorf("got %v", m["api_key"])
	}

	out, err := Serialize(sampleMap, TOML)
	if err != nil {
		t.Fatalf("Serialize TOML: %v", err)
	}
	if !strings.Contains(string(out), "${HOST}") {
		t.Errorf("expected placeholder in TOML output, got:\n%s", out)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	_, err := Parse([]byte("{}"), "xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestSerialize_InvalidFormat(t *testing.T) {
	_, err := Serialize(map[string]interface{}{}, "xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
