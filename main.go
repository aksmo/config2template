package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"config2template/converter"
)

func main() {
	input := flag.String("input", "", "Input JSON config file (required)")
	output := flag.String("output", "", "Output template file (default: <input>.tpl)")
	env := flag.String("env", "", "Output env file (default: <input>.env)")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: --input is required")
		flag.Usage()
		os.Exit(1)
	}

	raw, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", *input, err)
		os.Exit(1)
	}

	var config map[string]interface{}
	if err := json.Unmarshal(raw, &config); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	result := converter.Convert(config)

	// Write template file
	tplPath := *output
	if tplPath == "" {
		tplPath = *input + ".tpl"
	}
	tplData, err := json.MarshalIndent(result.Template, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling template: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(tplPath, append(tplData, '\n'), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing template: %v\n", err)
		os.Exit(1)
	}

	// Write env file
	envPath := *env
	if envPath == "" {
		envPath = *input + ".env"
	}
	keys := make([]string, 0, len(result.EnvVars))
	for k := range result.EnvVars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var lines []string
	for _, k := range keys {
		lines = append(lines, fmt.Sprintf("%s=%s", k, result.EnvVars[k]))
	}
	if err := os.WriteFile(envPath, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing env file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("template → %s\n", tplPath)
	fmt.Printf("env vars → %s\n", envPath)
}
