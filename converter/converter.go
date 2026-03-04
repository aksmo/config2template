// Package converter transforms JSON config files into templates and env var files.
// String leaf values are replaced with ${KEY_NAME} placeholders; numbers,
// booleans, and nulls are left as-is.
package converter

import (
	"fmt"
	"strings"
)

// Result holds the processed template tree and the extracted env vars.
type Result struct {
	Template map[string]interface{}
	EnvVars  map[string]string
}

// Convert walks a parsed JSON object and returns a Result.
func Convert(data map[string]interface{}) Result {
	envVars := make(map[string]string)
	template := processObject(data, nil, envVars)
	return Result{
		Template: template,
		EnvVars:  envVars,
	}
}

func processObject(obj map[string]interface{}, path []string, envVars map[string]string) map[string]interface{} {
	result := make(map[string]interface{}, len(obj))
	for k, v := range obj {
		result[k] = processValue(v, append(path, k), envVars)
	}
	return result
}

func processArray(arr []interface{}, path []string, envVars map[string]string) []interface{} {
	result := make([]interface{}, len(arr))
	for i, v := range arr {
		result[i] = processValue(v, append(path, fmt.Sprintf("%d", i)), envVars)
	}
	return result
}

func processValue(val interface{}, path []string, envVars map[string]string) interface{} {
	switch v := val.(type) {
	case string:
		key := pathToKey(path)
		envVars[key] = v
		return fmt.Sprintf("${%s}", key)
	case map[string]interface{}:
		return processObject(v, path, envVars)
	case []interface{}:
		return processArray(v, path, envVars)
	default:
		// numbers, booleans, null — keep as-is
		return val
	}
}

// pathToKey joins path segments with "_" and uppercases them,
// replacing non-alphanumeric characters with underscores.
func pathToKey(path []string) string {
	parts := make([]string, len(path))
	for i, p := range path {
		parts[i] = sanitize(p)
	}
	return strings.Join(parts, "_")
}

func sanitize(s string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(s) {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	return b.String()
}
