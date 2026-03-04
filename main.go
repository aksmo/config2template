package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"config2template/converter"
	"config2template/format"
)

func main() {
	input := flag.String("input", "", "Input config file — .json, .yaml/.yml, .toml (required)")
	output := flag.String("output", "", "Output template file (default: <input>.tpl)")
	env := flag.String("env", "", "Output env file (default: <input>.env)")
	fmtFlag := flag.String("format", "", "Force input format: json, yaml, toml (default: auto-detect from extension)")
	flag.Parse()

	if *input == "" {
		fmt.Fprintln(os.Stderr, "error: --input is required")
		flag.Usage()
		os.Exit(1)
	}

	// Determine format
	var fmt_ format.Format
	if *fmtFlag != "" {
		fmt_ = format.Format(*fmtFlag)
	} else {
		var err error
		fmt_, err = format.Detect(*input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	// Read and parse input
	raw, err := os.ReadFile(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", *input, err)
		os.Exit(1)
	}
	config, err := format.Parse(raw, fmt_)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %s: %v\n", fmt_, err)
		os.Exit(1)
	}

	result := converter.Convert(config)

	// Write template file
	tplPath := *output
	if tplPath == "" {
		tplPath = *input + ".tpl"
	}
	tplData, err := format.Serialize(result.Template, fmt_)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error serializing template: %v\n", err)
		os.Exit(1)
	}
	if err := os.WriteFile(tplPath, tplData, 0644); err != nil {
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
