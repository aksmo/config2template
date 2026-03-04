package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"config2template/converter"
	"config2template/format"
)

func main() {
	input := flag.String("input", "", "Input config file — .json, .yaml/.yml, .toml, .env (required, or omit for interactive mode)")
	output := flag.String("output", "", "Output template file (default: <input>.tpl)")
	env := flag.String("env", "", "Output env file (default: <input>.env)")
	fmtFlag := flag.String("format", "", "Force input format: json, yaml, toml, env (default: auto-detect from extension)")
	flag.Parse()

	if *input == "" {
		if isTerminal() {
			interactive()
			return
		}
		fmt.Fprintln(os.Stderr, "error: --input is required")
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*input, *output, *env, *fmtFlag); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run performs the conversion given resolved paths and an optional format override.
func run(inputPath, outputPath, envPath, fmtOverride string) error {
	// Determine format
	var fmt_ format.Format
	if fmtOverride != "" {
		fmt_ = format.Format(fmtOverride)
	} else {
		var err error
		fmt_, err = format.Detect(inputPath)
		if err != nil {
			return err
		}
	}

	raw, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading %s: %w", inputPath, err)
	}
	config, err := format.Parse(raw, fmt_)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", fmt_, err)
	}

	result := converter.Convert(config)

	// Write template
	tplPath := outputPath
	if tplPath == "" {
		tplPath = inputPath + ".tpl"
	}
	tplData, err := format.Serialize(result.Template, fmt_)
	if err != nil {
		return fmt.Errorf("serializing template: %w", err)
	}
	if err := os.WriteFile(tplPath, tplData, 0644); err != nil {
		return fmt.Errorf("writing template: %w", err)
	}

	// Write env file
	ep := envPath
	if ep == "" {
		ep = inputPath + ".env"
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
	if err := os.WriteFile(ep, []byte(strings.Join(lines, "\n")+"\n"), 0644); err != nil {
		return fmt.Errorf("writing env file: %w", err)
	}

	fmt.Printf("template → %s\n", tplPath)
	fmt.Printf("env vars → %s  (%d variables)\n", ep, len(result.EnvVars))
	return nil
}

// interactive runs a simple prompt-driven wizard.
func interactive() {
	r := bufio.NewReader(os.Stdin)

	fmt.Println("config2template — interactive mode")
	fmt.Println("Convert a config file into a template + env file.")
	fmt.Println()

	inputPath := prompt(r, "Input file: ", "")
	if inputPath == "" {
		fmt.Fprintln(os.Stderr, "error: input file is required")
		os.Exit(1)
	}

	defaultTpl := inputPath + ".tpl"
	defaultEnv := inputPath + ".env"

	// Try to detect format; if unknown, ask.
	fmt_, err := format.Detect(inputPath)
	if err != nil {
		fmtStr := prompt(r, "Format (json/yaml/toml/env): ", "")
		fmt_ = format.Format(strings.TrimSpace(fmtStr))
	} else {
		fmt.Printf("Format detected: %s\n", fmt_)
	}

	tplPath := prompt(r, fmt.Sprintf("Output template [%s]: ", defaultTpl), defaultTpl)
	envPath := prompt(r, fmt.Sprintf("Output env file [%s]: ", defaultEnv), defaultEnv)

	fmt.Println()
	if err := run(inputPath, tplPath, envPath, string(fmt_)); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// prompt prints a message and reads a line from r.
// If the user enters nothing, it returns the defaultVal.
func prompt(r *bufio.Reader, message, defaultVal string) string {
	fmt.Print(message)
	line, _ := r.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}

// isTerminal reports whether stdin is an interactive terminal.
func isTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}
