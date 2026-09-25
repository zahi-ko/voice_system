package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"voice_system/internal/adapters/audioio"
	"voice_system/internal/adapters/memory"
	"voice_system/internal/application"
	"voice_system/internal/transport/http/generated"

	"github.com/labstack/echo/v5"
)

func TestAudioEndpointsSupportFixtureFormats(t *testing.T) {
	fixtures, err := filepath.Glob(filepath.Join("..", "..", "..", "..", "tests", "data", "test.*"))
	if err != nil {
		t.Fatal(err)
	}
	if len(fixtures) == 0 {
		t.Fatal("no audio fixtures found")
	}

	for _, fixture := range fixtures {
		t.Run(filepath.Ext(fixture), func(t *testing.T) {
			content, err := os.ReadFile(fixture)
			if err != nil {
				t.Fatal(err)
			}

			store, err := memory.NewAudioStore(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			handler := NewHandler(application.NewAudioService(audioio.NewDecoder(), store))
			router := echo.New()
			Register(router, handler)

			uploadBody, contentType := multipartFixture(t, filepath.Base(fixture), content)
			uploadRequest := httptest.NewRequest(http.MethodPost, "/audio/upload", uploadBody)
			uploadRequest.Header.Set(echo.HeaderContentType, contentType)
			uploadResponse := httptest.NewRecorder()
			router.ServeHTTP(uploadResponse, uploadRequest)

			if uploadResponse.Code != http.StatusCreated {
				t.Fatalf("upload status = %d, want %d: %s", uploadResponse.Code, http.StatusCreated, uploadResponse.Body.String())
			}

			var uploaded generated.AudioMeta
			if err := json.Unmarshal(uploadResponse.Body.Bytes(), &uploaded); err != nil {
				t.Fatalf("decode upload response: %v", err)
			}
			if uploaded.ID == nil || *uploaded.ID == "" {
				t.Fatal("upload response has no audio ID")
			}
			if uploaded.Name == nil || *uploaded.Name != filepath.Base(fixture) {
				t.Fatalf("upload response name = %v, want %s", uploaded.Name, filepath.Base(fixture))
			}

			listRequest := httptest.NewRequest(http.MethodGet, "/audio", nil)
			listResponse := httptest.NewRecorder()
			router.ServeHTTP(listResponse, listRequest)
			if listResponse.Code != http.StatusOK {
				t.Fatalf("list status = %d, want %d", listResponse.Code, http.StatusOK)
			}
			var listed generated.AudioList
			if err := json.Unmarshal(listResponse.Body.Bytes(), &listed); err != nil {
				t.Fatalf("decode list response: %v", err)
			}
			if len(listed.Items) != 1 || listed.Items[0].ID == nil || *listed.Items[0].ID != *uploaded.ID {
				t.Fatalf("list response = %+v, want uploaded audio", listed)
			}

			downloadRequest := httptest.NewRequest(http.MethodGet, "/audio/"+*uploaded.ID+"/download", nil)
			downloadResponse := httptest.NewRecorder()
			router.ServeHTTP(downloadResponse, downloadRequest)
			if downloadResponse.Code != http.StatusOK {
				t.Fatalf("download status = %d, want %d", downloadResponse.Code, http.StatusOK)
			}
			if !bytes.Equal(downloadResponse.Body.Bytes(), content) {
				t.Fatal("downloaded audio does not match uploaded fixture")
			}

			deleteRequest := httptest.NewRequest(http.MethodDelete, "/audio/"+*uploaded.ID, nil)
			deleteResponse := httptest.NewRecorder()
			router.ServeHTTP(deleteResponse, deleteRequest)
			if deleteResponse.Code != http.StatusNoContent {
				t.Fatalf("delete status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
			}
		})
	}
}

func multipartFixture(t *testing.T, name string, content []byte) (*bytes.Reader, string) {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	return bytes.NewReader(body.Bytes()), writer.FormDataContentType()
}
