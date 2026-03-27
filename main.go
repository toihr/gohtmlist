package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed template.html
var templates embed.FS

var tmpl = template.Must(template.ParseFS(templates, "template.html"))

type PageData struct{
	CSV string
}

var csvContent = `
		name,age
		Alice,30
		Bob,25
`
func indentCSV(csv string, indent string) string {
    lines := strings.Split(csv, "\n")
    for i, l := range lines {
        lines[i] = indent + l
    }
    return strings.Join(lines, "\n")
}

func Render(csvPath, outputPath string) error {
	csvBytes, err := os.ReadFile(csvPath)
	if err != nil {
		return err
	}

    f, err := os.Create(outputPath)
    if err != nil {
        return err
    }
    defer f.Close()

    return tmpl.Execute(f, PageData{CSV: indentCSV(string(csvBytes), "        ")})
}



func runTests() bool {
	dir, err := os.MkdirTemp("", "gohtmlist-test-*")
	if err != nil {
		fmt.Printf("failed to create temp dir: %v\n", err)
		return false
	}
	defer os.RemoveAll(dir)

	csvPath := filepath.Join(dir, "input.csv")
	outPath := filepath.Join(dir, "output.html")

	if err := os.WriteFile(csvPath, []byte("name,age\nAlice,30\nBob,25\n"), 0600); err != nil {
		fmt.Printf("failed to write sample CSV: %v\n", err)
		return false
	}

	if err := Render(csvPath, outPath); err != nil {
		fmt.Printf("Render failed: %v\n", err)
		return false
	}

	htmlBytes, err := os.ReadFile(outPath)
	if err != nil {
		fmt.Printf("failed to read output: %v\n", err)
		return false
	}
	html := string(htmlBytes)

	type testCase struct {
		name string
		fn   func() error
	}

	tests := []testCase{
		{
			name: "OutputFileCreated",
			fn: func() error {
				if _, err := os.Stat(outPath); os.IsNotExist(err) {
					return fmt.Errorf("output file was not created")
				}
				return nil
			},
		},
		{
			name: "OutputContainsCSVData",
			fn: func() error {
				for _, want := range []string{"name,age", "Alice,30", "Bob,25"} {
					if !strings.Contains(html, want) {
						return fmt.Errorf("output HTML does not contain %q", want)
					}
				}
				return nil
			},
		},
		{
			name: "OutputIsValidHTML",
			fn: func() error {
				lowerHTML := strings.ToLower(html)
				for _, tag := range []string{"<!doctype html>", "<html", "</html>", "<head>", "</head>", "<body>", "</body>"} {
					if !strings.Contains(lowerHTML, tag) {
						return fmt.Errorf("output HTML missing expected tag %q", tag)
					}
				}
				return nil
			},
		},
		{
			name: "OutputContainsCSVScriptTag",
			fn: func() error {
				if !strings.Contains(html, `type="text/csv"`) {
					return fmt.Errorf(`output HTML missing <script type="text/csv"> tag`)
				}
				return nil
			},
		},
	}

	passed, failed := 0, 0
	for _, tc := range tests {
		if err := tc.fn(); err != nil {
			fmt.Printf("--- FAIL: %s\n    %v\n", tc.name, err)
			failed++
		} else {
			fmt.Printf("--- PASS: %s\n", tc.name)
			passed++
		}
	}

	fmt.Printf("\n%d passed, %d failed\n", passed, failed)
	return failed == 0
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "test" {
		if !runTests() {
			os.Exit(1)
		}
		return
	}

	filePath := flag.String("f", "", "Path to input CSV file")
	outPath := flag.String("out", "output.html", "Path to output HTML file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: gohtmlist [command] [flags]\n\nCommands:\n  test\tRun built-in self-tests\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *filePath == "" {
		log.Panic("Input CSV file path is required. Use -f flag to specify it.")
	}

	if err := Render(*filePath, *outPath); err != nil {
		log.Panic(err)
	}

	log.Printf("HTML file generated at: %s", *outPath)
}
