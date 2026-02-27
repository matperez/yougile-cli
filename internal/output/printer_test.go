package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestPrintTable_HeadersAndRows_WritesTable(t *testing.T) {
	var buf bytes.Buffer
	headers := []string{"ID", "Name"}
	rows := [][]string{
		{"1", "Alice"},
		{"2", "Bob"},
	}
	err := PrintTable(&buf, headers, rows)
	if err != nil {
		t.Fatalf("PrintTable: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "ID") || !strings.Contains(out, "Name") {
		t.Errorf("output missing headers: %q", out)
	}
	if !strings.Contains(out, "Alice") || !strings.Contains(out, "Bob") {
		t.Errorf("output missing row data: %q", out)
	}
}

func TestPrintJSON_ValidStruct_WritesValidJSON(t *testing.T) {
	var buf bytes.Buffer
	v := map[string]string{"a": "b", "c": "d"}
	err := PrintJSON(&buf, v)
	if err != nil {
		t.Fatalf("PrintJSON: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Errorf("output is not valid JSON: %q", buf.String())
	}
	if decoded["a"] != "b" || decoded["c"] != "d" {
		t.Errorf("decoded = %v", decoded)
	}
}
