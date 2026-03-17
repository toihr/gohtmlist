# gohtmlist

**gohtmlist** is a Go CLI tool that converts CSV files into interactive, self-contained HTML tables. The generated HTML file requires no external dependencies — everything (styles, JavaScript, and data) is embedded directly in a single file.

## Features

- **Interactive filtering** — a search box filters all rows in real time
- **Multi-column sorting** — click any column header to sort ascending or descending
- **Search highlighting** — matching text is highlighted in the results
- **Row counter** — shows the number of visible rows after filtering
- **Virtual scrolling** — only visible rows are rendered, keeping the page fast even for large datasets
- **Zero runtime dependencies** — the output is a single, portable `.html` file

## Installation

Requires [Go 1.21+](https://golang.org/dl/).

```bash
go install github.com/toihr/gohtmlist@latest
```

Or clone and build from source:

```bash
git clone https://github.com/toihr/gohtmlist.git
cd gohtmlist
go build -o gohtmlist .
```

## Usage

```
gohtmlist -f <input.csv> [-out <output.html>]
```

| Flag   | Default         | Description                        |
|--------|-----------------|------------------------------------|
| `-f`   | *(required)*    | Path to the input CSV file         |
| `-out` | `output.html`   | Path for the generated HTML file   |

### Example

```bash
gohtmlist -f data.csv -out report.html
```

This reads `data.csv` and writes an interactive HTML table to `report.html`.

### Sample CSV

```csv
name,age,city,salary
Alice,30,New York,85000
Bob,25,Los Angeles,72000
Carol,35,Chicago,91000
```

Open the generated `report.html` in any modern browser — no server required.

## Go API

`gohtmlist` also exposes a `Render` function for use as a library:

```go
import "github.com/toihr/gohtmlist"

err := gohtmlist.Render("data.csv", "output.html")
if err != nil {
    log.Fatal(err)
}
```

### `Render(csvPath, outputPath string) error`

Reads the CSV file at `csvPath`, processes it through the embedded HTML template, and writes the result to `outputPath`. Returns an error if either file operation fails.

## Development

```bash
# Run from source
go run . -f test/sample.csv -out test/output.html

# Build a binary
go build -o gohtmlist .
```

## License

This project is open source. See the repository for license details.
