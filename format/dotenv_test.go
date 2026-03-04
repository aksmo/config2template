package format

import (
	"strings"
	"testing"
)

func TestParseDotenv_Basic(t *testing.T) {
	input := []byte("DB_HOST=localhost\nDB_PASS=s3cr3t\nPORT=5432\n")
	m, err := parseDotenv(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["DB_HOST"] != "localhost" {
		t.Errorf("got %v", m["DB_HOST"])
	}
	if m["DB_PASS"] != "s3cr3t" {
		t.Errorf("got %v", m["DB_PASS"])
	}
	if m["PORT"] != "5432" {
		t.Errorf("got %v", m["PORT"])
	}
}

func TestParseDotenv_Quotes(t *testing.T) {
	input := []byte(`
DB_URL="postgresql://user:pass@host/db"
GREETING='hello world'
`)
	m, err := parseDotenv(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["DB_URL"] != "postgresql://user:pass@host/db" {
		t.Errorf("got %v", m["DB_URL"])
	}
	if m["GREETING"] != "hello world" {
		t.Errorf("got %v", m["GREETING"])
	}
}

func TestParseDotenv_CommentsAndBlanks(t *testing.T) {
	input := []byte(`
# this is a comment
KEY=value # inline comment

ANOTHER=thing
`)
	m, err := parseDotenv(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["KEY"] != "value" {
		t.Errorf("got %q", m["KEY"])
	}
	if m["ANOTHER"] != "thing" {
		t.Errorf("got %q", m["ANOTHER"])
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
}

func TestParseDotenv_Export(t *testing.T) {
	input := []byte("export API_KEY=abc123\n")
	m, err := parseDotenv(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m["API_KEY"] != "abc123" {
		t.Errorf("got %v", m["API_KEY"])
	}
}

func TestParseDotenv_InvalidLine(t *testing.T) {
	_, err := parseDotenv([]byte("NOEQUALSSIGN\n"))
	if err == nil {
		t.Error("expected error for line without '='")
	}
}

func TestSerializeDotenv_Basic(t *testing.T) {
	m := map[string]interface{}{
		"API_KEY": "${API_KEY}",
		"DB_HOST": "${DB_HOST}",
	}
	out, err := serializeDotenv(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := string(out)
	if !strings.Contains(s, "API_KEY=${API_KEY}") {
		t.Errorf("missing API_KEY line in:\n%s", s)
	}
	if !strings.Contains(s, "DB_HOST=${DB_HOST}") {
		t.Errorf("missing DB_HOST line in:\n%s", s)
	}
}

func TestSerializeDotenv_QuotesSpaces(t *testing.T) {
	m := map[string]interface{}{
		"GREETING": "hello world",
	}
	out, err := serializeDotenv(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(out), `GREETING="hello world"`) {
		t.Errorf("expected quoted value, got:\n%s", out)
	}
}

func TestDetect_DotEnv(t *testing.T) {
	f, err := Detect(".env")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f != DotEnv {
		t.Errorf("expected DotEnv, got %v", f)
	}
}

func TestParseAndSerialize_DotEnv(t *testing.T) {
	input := []byte("API_KEY=secret\nDB_HOST=localhost\n")
	m, err := Parse(input, DotEnv)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if m["API_KEY"] != "secret" {
		t.Errorf("got %v", m["API_KEY"])
	}

	out, err := Serialize(map[string]interface{}{"API_KEY": "${API_KEY}"}, DotEnv)
	if err != nil {
		t.Fatalf("Serialize: %v", err)
	}
	if !strings.Contains(string(out), "API_KEY=${API_KEY}") {
		t.Errorf("got:\n%s", out)
	}
}
