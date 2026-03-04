package converter

import (
	"testing"
)

func TestConvert_FlatStrings(t *testing.T) {
	input := map[string]interface{}{
		"host": "localhost",
		"port": float64(5432), // JSON numbers unmarshal as float64
	}
	result := Convert(input)

	if result.Template["host"] != "${HOST}" {
		t.Errorf("expected ${HOST}, got %v", result.Template["host"])
	}
	if result.Template["port"] != float64(5432) {
		t.Errorf("expected port 5432 unchanged, got %v", result.Template["port"])
	}
	if result.EnvVars["HOST"] != "localhost" {
		t.Errorf("expected HOST=localhost, got %v", result.EnvVars["HOST"])
	}
	if _, ok := result.EnvVars["PORT"]; ok {
		t.Error("numeric port should not appear in env vars")
	}
}

func TestConvert_NestedObject(t *testing.T) {
	input := map[string]interface{}{
		"database": map[string]interface{}{
			"password": "s3cr3t",
			"name":     "mydb",
		},
	}
	result := Convert(input)

	db, ok := result.Template["database"].(map[string]interface{})
	if !ok {
		t.Fatal("expected database to be a nested object in template")
	}
	if db["password"] != "${DATABASE_PASSWORD}" {
		t.Errorf("got %v", db["password"])
	}
	if db["name"] != "${DATABASE_NAME}" {
		t.Errorf("got %v", db["name"])
	}
	if result.EnvVars["DATABASE_PASSWORD"] != "s3cr3t" {
		t.Errorf("got %v", result.EnvVars["DATABASE_PASSWORD"])
	}
}

func TestConvert_Array(t *testing.T) {
	input := map[string]interface{}{
		"hosts": []interface{}{"host1.example.com", "host2.example.com"},
	}
	result := Convert(input)

	arr, ok := result.Template["hosts"].([]interface{})
	if !ok {
		t.Fatal("expected hosts to be an array in template")
	}
	if arr[0] != "${HOSTS_0}" {
		t.Errorf("got %v", arr[0])
	}
	if arr[1] != "${HOSTS_1}" {
		t.Errorf("got %v", arr[1])
	}
	if result.EnvVars["HOSTS_0"] != "host1.example.com" {
		t.Errorf("got %v", result.EnvVars["HOSTS_0"])
	}
}

func TestConvert_BoolAndNull(t *testing.T) {
	input := map[string]interface{}{
		"enabled": true,
		"nothing": nil,
	}
	result := Convert(input)

	if result.Template["enabled"] != true {
		t.Errorf("boolean should pass through unchanged")
	}
	if result.Template["nothing"] != nil {
		t.Errorf("null should pass through unchanged")
	}
	if len(result.EnvVars) != 0 {
		t.Errorf("no env vars expected for bool/null, got %v", result.EnvVars)
	}
}

func TestConvert_SpecialCharsInKey(t *testing.T) {
	input := map[string]interface{}{
		"my-key": "value",
	}
	result := Convert(input)

	if result.Template["my-key"] != "${MY_KEY}" {
		t.Errorf("got %v", result.Template["my-key"])
	}
	if result.EnvVars["MY_KEY"] != "value" {
		t.Errorf("got %v", result.EnvVars["MY_KEY"])
	}
}

func TestPathToKey(t *testing.T) {
	cases := []struct {
		path []string
		want string
	}{
		{[]string{"host"}, "HOST"},
		{[]string{"database", "password"}, "DATABASE_PASSWORD"},
		{[]string{"my-key"}, "MY_KEY"},
		{[]string{"api.key"}, "API_KEY"},
		{[]string{"servers", "0"}, "SERVERS_0"},
	}
	for _, c := range cases {
		got := pathToKey(c.path)
		if got != c.want {
			t.Errorf("pathToKey(%v) = %q, want %q", c.path, got, c.want)
		}
	}
}
