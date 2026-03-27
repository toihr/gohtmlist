package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIndentCSV_NoIndent(t *testing.T) {
	input := "a,b\n1,2"
	got := indentCSV(input, "")
	if got != input {
		t.Errorf("indentCSV with empty indent changed content: got %q, want %q", got, input)
	}
}

func TestIndentCSV_WithIndent(t *testing.T) {
	input := "a,b\n1,2"
	indent := "    "
	got := indentCSV(input, indent)
	for _, line := range strings.Split(got, "\n") {
		if !strings.HasPrefix(line, indent) {
			t.Errorf("line %q does not start with indent %q", line, indent)
		}
	}
}

func TestIndentCSV_PreservesLineCount(t *testing.T) {
	input := "h1,h2\nv1,v2\nv3,v4"
	got := indentCSV(input, "\t")
	if strings.Count(got, "\n") != strings.Count(input, "\n") {
		t.Errorf("indentCSV changed number of newlines: got %d, want %d",
			strings.Count(got, "\n"), strings.Count(input, "\n"))
	}
}

func TestRender_MissingInputFile(t *testing.T) {
	err := Render("/nonexistent/path/file.csv", filepath.Join(t.TempDir(), "out.html"))
	if err == nil {
		t.Error("expected error for missing input file, got nil")
	}
}

func TestRender_CreatesOutputFile(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "input.csv")
	outPath := filepath.Join(dir, "output.html")

	if err := os.WriteFile(csvPath, []byte("name,age\nAlice,30\nBob,25\n"), 0644); err != nil {
		t.Fatalf("failed to write temp CSV: %v", err)
	}

	if err := Render(csvPath, outPath); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestRender_OutputContainsCSVData(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "input.csv")
	outPath := filepath.Join(dir, "output.html")

	csvContent := "name,age\nAlice,30\nBob,25\n"
	if err := os.WriteFile(csvPath, []byte(csvContent), 0644); err != nil {
		t.Fatalf("failed to write temp CSV: %v", err)
	}

	if err := Render(csvPath, outPath); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	html := string(data)

	for _, want := range []string{"name,age", "Alice,30", "Bob,25"} {
		if !strings.Contains(html, want) {
			t.Errorf("output HTML does not contain %q", want)
		}
	}
}

func TestRender_OutputIsValidHTML(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "input.csv")
	outPath := filepath.Join(dir, "output.html")

	if err := os.WriteFile(csvPath, []byte("col1,col2\nval1,val2\n"), 0644); err != nil {
		t.Fatalf("failed to write temp CSV: %v", err)
	}

	if err := Render(csvPath, outPath); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	html := string(data)

	for _, tag := range []string{"<!doctype html>", "<html", "</html>", "<head>", "</head>", "<body>", "</body>"} {
		if !strings.Contains(strings.ToLower(html), strings.ToLower(tag)) {
			t.Errorf("output HTML missing expected tag %q", tag)
		}
	}
}

func TestRunTests_AllPass(t *testing.T) {
	if !runTests() {
		t.Error("runTests() returned false; expected all built-in tests to pass")
	}
}

func TestRender_OutputContainsCSVScriptTag(t *testing.T) {
	dir := t.TempDir()
	csvPath := filepath.Join(dir, "input.csv")
	outPath := filepath.Join(dir, "output.html")

	if err := os.WriteFile(csvPath, []byte("x,y\n1,2\n"), 0644); err != nil {
		t.Fatalf("failed to write temp CSV: %v", err)
	}

	if err := Render(csvPath, outPath); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}
	html := string(data)

	if !strings.Contains(html, `type="text/csv"`) {
		t.Error(`output HTML missing <script type="text/csv"> tag`)
	}
}
