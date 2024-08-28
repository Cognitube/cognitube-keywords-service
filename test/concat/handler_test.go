package concat

import (
	"bytes"
	"cognitube.com/keywords-service/server/handler"
	"cognitube.com/keywords-service/server/service/concat"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConcatJSONByOrder(t *testing.T) {
	// Create a buffer to hold the multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Create sample JSON files
	json1 := `{"keyword": "xxx", "description": "xxx", "start": "0:12", "end": "1:24", "clip": 0}`
	json2 := `{"keyword": "xxx", "description": "xxx", "start": "5:10", "end": "6:02", "clip": 0}`
	json3 := `{"keyword": "xxx", "description": "xxx", "start": "0:16", "end": "1:24", "clip": 1}`

	// Add JSON files to the multipart form
	part1, _ := writer.CreateFormFile("file", "file1.json")
	part1.Write([]byte(json1))
	part2, _ := writer.CreateFormFile("file", "file2.json")
	part2.Write([]byte(json2))
	part3, _ := writer.CreateFormFile("file", "file3.json")
	part3.Write([]byte(json3))

	writer.Close()

	// Create a new HTTP request with the multipart form data
	req := httptest.NewRequest("POST", "/api/v1/concat", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create a new HTTP response recorder
	rr := httptest.NewRecorder()

	// Create a new ConcatHandler
	concatService := concat.NewConcatService()
	h := handler.NewConcatHandler(concatService)

	// Call the ConcatJSONByOrder handler
	h.ConcatJSONByOrder(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
	expected := []concat.ConcatJSON{
		{Keyword: "xxx", Description: "xxx", Start: "0:12", End: "1:24", Clip: 0},
		{Keyword: "xxx", Description: "xxx", Start: "5:10", End: "6:02", Clip: 0},
		{Keyword: "xxx", Description: "xxx", Start: "0:16", End: "1:24", Clip: 1},
	}
	var actual []concat.ConcatJSON
	if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
		t.Errorf("error decoding response body: %v", err)
	}

	if !equal(expected, actual) {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func equal(a, b []concat.ConcatJSON) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
