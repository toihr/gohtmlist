package main

import (
	"embed"
	"flag"
	"html/template"
	"log"
	"os"
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



func main() {
	filePath := flag.String("f", "", "Path to input CSV file")
	outPath := flag.String("out", "output.html", "Path to output HTML file")
	flag.Parse()

	if *filePath == "" {
		log.Panic("Input CSV file path is required. Use -f flag to specify it.")
	}


	if err := Render(*filePath,*outPath); err != nil {
		log.Panic(err)
	}


	log.Printf("HTML file generated at: %s", *outPath)
}
