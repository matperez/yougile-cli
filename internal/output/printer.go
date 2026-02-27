package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// PrintTable writes a table with headers and rows to w.
// Each row must have the same number of columns as headers.
func PrintTable(w io.Writer, headers []string, rows [][]string) error {
	if len(headers) == 0 {
		return nil
	}
	widths := make([]int, len(headers))
	for i, h := range headers {
		if len(h) > widths[i] {
			widths[i] = len(h)
		}
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}
	format := strings.TrimSuffix(strings.Repeat("%-*s  ", len(headers)), "  ") + "\n"
	args := make([]interface{}, 0, len(headers)*2)
	for i := range headers {
		args = append(args, widths[i], headers[i])
	}
	if _, err := fmt.Fprintf(w, format, args...); err != nil {
		return err
	}
	sep := strings.Repeat("-", sum(widths)+2*len(headers)-2) + "\n"
	if _, err := io.WriteString(w, sep); err != nil {
		return err
	}
	for _, row := range rows {
		args = args[:0]
		for i := 0; i < len(headers) && i < len(row); i++ {
			args = append(args, widths[i], row[i])
		}
		if _, err := fmt.Fprintf(w, format, args...); err != nil {
			return err
		}
	}
	return nil
}

func sum(w []int) int {
	s := 0
	for _, x := range w {
		s += x
	}
	return s
}

// PrintJSON writes v as JSON to w. No indentation.
func PrintJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}
