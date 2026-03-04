package format

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// parseDotenv parses a .env file into a map.
// Supports:
//   - KEY=value
//   - KEY="value" or KEY='value'  (quotes stripped)
//   - export KEY=value
//   - # comment lines and inline comments
//   - blank lines
func parseDotenv(data []byte) (map[string]interface{}, error) {
	out := make(map[string]interface{})
	scanner := bufio.NewScanner(bytes.NewReader(data))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// strip optional leading "export "
		line = strings.TrimPrefix(line, "export ")

		eq := strings.IndexByte(line, '=')
		if eq < 0 {
			return nil, fmt.Errorf("line %d: missing '=' in %q", lineNum, line)
		}
		key := strings.TrimSpace(line[:eq])
		val := line[eq+1:]

		// strip inline comment (outside of quotes)
		val = stripInlineComment(val)
		val = unquote(val)

		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", lineNum)
		}
		out[key] = val
	}
	return out, scanner.Err()
}

// serializeDotenv writes a map as a .env file (keys sorted).
func serializeDotenv(data map[string]interface{}) ([]byte, error) {
	keys := make([]string, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for _, k := range keys {
		v := data[k]
		switch val := v.(type) {
		case string:
			// quote if value contains spaces or special chars
			if needsQuoting(val) {
				fmt.Fprintf(&buf, "%s=\"%s\"\n", k, escapeDoubleQuotes(val))
			} else {
				fmt.Fprintf(&buf, "%s=%s\n", k, val)
			}
		default:
			fmt.Fprintf(&buf, "%s=%v\n", k, val)
		}
	}
	return buf.Bytes(), nil
}

func stripInlineComment(s string) string {
	inSingle, inDouble := false, false
	for i, r := range s {
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
			}
		case '#':
			if !inSingle && !inDouble {
				return strings.TrimSpace(s[:i])
			}
		}
	}
	return s
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') ||
			(s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}

func needsQuoting(s string) bool {
	return strings.ContainsAny(s, " \t#\"'\\")
}

func escapeDoubleQuotes(s string) string {
	return strings.ReplaceAll(s, `"`, `\"`)
}
