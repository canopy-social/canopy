package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/sumi-devs/canopy-social/canopy/internal/auth"
	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
)

type mockRepository struct {
	db map[string]*Attachment
}

func (m *mockRepository) Create(ctx context.Context, att *Attachment) (*Attachment, error) {
	m.db[att.ID] = att
	return att, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id string) (*Attachment, error) {
	att, ok := m.db[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return att, nil
}

func (m *mockRepository) Update(ctx context.Context, id string, altText *string) (*Attachment, error) {
	att, ok := m.db[id]
	if !ok {
		return nil, errors.New("not found")
	}
	att.AltText = altText
	return att, nil
}

func (m *mockRepository) Delete(ctx context.Context, id string) error {
	delete(m.db, id)
	return nil
}

type mockS3Client struct {
	uploads map[string][]byte
	deletes []string
}

func (m *mockS3Client) Upload(ctx context.Context, key string, contentType string, reader io.Reader) (string, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	m.uploads[key] = data
	return "http://localhost:9000/canopy-media/" + key, nil
}

func (m *mockS3Client) Delete(ctx context.Context, key string) error {
	m.deletes = append(m.deletes, key)
	delete(m.uploads, key)
	return nil
}

func createTestImage(w, h int) []byte {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes()
}

func TestProcessUpload_Image(t *testing.T) {
	repo := &mockRepository{db: make(map[string]*Attachment)}
	s3 := &mockS3Client{uploads: make(map[string][]byte)}
	cfg := &config.Config{}
	cfg.Features.MaxImageMB = 10

	svc := NewService(repo, s3, cfg)

	imgBytes := createTestImage(100, 100)
	att, err := svc.ProcessUpload(context.Background(), "user123", bytes.NewReader(imgBytes), "avatar.png")
	if err != nil {
		t.Fatalf("failed to process image upload: %v", err)
	}

	if att.Type != "image" {
		t.Errorf("expected type to be image, got %s", att.Type)
	}

	if *att.Width != 100 || *att.Height != 100 {
		t.Errorf("expected dimensions 100x100, got %dx%d", *att.Width, *att.Height)
	}

	if att.Blurhash == nil || *att.Blurhash == "" {
		t.Error("expected blurhash to be generated")
	}

	if att.ThumbnailURL == nil || *att.ThumbnailURL == "" {
		t.Error("expected thumbnail URL to be set")
	}

	if len(s3.uploads) != 2 {
		t.Errorf("expected 2 files uploaded to S3 (original + thumbnail), got %d", len(s3.uploads))
	}
}

func TestProcessUpload_Resizing(t *testing.T) {
	repo := &mockRepository{db: make(map[string]*Attachment)}
	s3 := &mockS3Client{uploads: make(map[string][]byte)}
	cfg := &config.Config{}
	cfg.Features.MaxImageMB = 10

	svc := NewService(repo, s3, cfg)

	imgBytes := createTestImage(600, 300)
	att, err := svc.ProcessUpload(context.Background(), "user123", bytes.NewReader(imgBytes), "wide.png")
	if err != nil {
		t.Fatalf("failed: %v", err)
	}

	if *att.Width != 600 || *att.Height != 300 {
		t.Errorf("expected original dimensions 600x300, got %dx%d", *att.Width, *att.Height)
	}

	hasThumb := false
	for k, data := range s3.uploads {
		if bytes.Contains([]byte(k), []byte("_thumb.jpg")) {
			hasThumb = true
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				t.Fatalf("failed to decode generated thumbnail: %v", err)
			}
			bounds := img.Bounds()
			if bounds.Dx() != 400 || bounds.Dy() != 200 {
				t.Errorf("expected thumbnail resized to 400x200 (aspect ratio), got %dx%d", bounds.Dx(), bounds.Dy())
			}
		}
	}

	if !hasThumb {
		t.Error("expected thumbnail file to be uploaded")
	}
}

func TestProcessUpload_Video(t *testing.T) {
	repo := &mockRepository{db: make(map[string]*Attachment)}
	s3 := &mockS3Client{uploads: make(map[string][]byte)}
	cfg := &config.Config{}
	cfg.Features.MaxVideoMB = 100

	svc := NewService(repo, s3, cfg)

	fakeVideo := []byte("\x00\x00\x00\x18ftypmp42\x00\x00\x00\x00mp42isomdummyvideodata")
	att, err := svc.ProcessUpload(context.Background(), "user123", bytes.NewReader(fakeVideo), "clip.mp4")
	if err != nil {
		t.Fatalf("failed to process video: %v", err)
	}

	if att.Type != "video" {
		t.Errorf("expected video type, got %s", att.Type)
	}

	if att.ThumbnailURL != nil {
		t.Error("expected no thumbnail for video")
	}

	if att.Blurhash != nil {
		t.Error("expected no blurhash for video")
	}

	if len(s3.uploads) != 1 {
		t.Errorf("expected only 1 file uploaded, got %d", len(s3.uploads))
	}
}

func TestProcessUpload_SizeValidation(t *testing.T) {
	repo := &mockRepository{db: make(map[string]*Attachment)}
	s3 := &mockS3Client{uploads: make(map[string][]byte)}
	cfg := &config.Config{}
	cfg.Features.MaxImageMB = 1

	svc := NewService(repo, s3, cfg)

	imgBytes := append([]byte("\x89PNG\r\n\x1a\n"), make([]byte, 2*1024*1024)...)
	_, err := svc.ProcessUpload(context.Background(), "user123", bytes.NewReader(imgBytes), "huge.png")
	if err == nil {
		t.Error("expected size validation error, got nil")
	}
}

func TestHandler_CRUD(t *testing.T) {
	repo := &mockRepository{db: make(map[string]*Attachment)}
	s3 := &mockS3Client{uploads: make(map[string][]byte)}
	cfg := &config.Config{}
	cfg.Features.MaxImageMB = 10

	svc := NewService(repo, s3, cfg)
	h := NewHandler(svc)

	jwtMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), auth.ContextKeyAccountID, "user123")
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	r := chi.NewRouter()
	h.RegisterRoutes(r, jwtMiddleware)

	server := httptest.NewServer(r)
	defer server.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, _ := writer.CreateFormFile("file", "test.png")
	part.Write(createTestImage(50, 50))
	writer.Close()

	req, _ := http.NewRequest("POST", server.URL+"/api/v1/media", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	var att Attachment
	json.NewDecoder(resp.Body).Decode(&att)
	if att.ID == "" {
		t.Error("expected valid created attachment ID")
	}

	req2, _ := http.NewRequest("GET", server.URL+"/api/v1/media/"+att.ID, nil)
	resp2, _ := http.DefaultClient.Do(req2)
	if resp2.StatusCode != http.StatusOK {
		t.Errorf("GET failed: %d", resp2.StatusCode)
	}

	updateData := map[string]string{"alt_text": "Beautiful landscape"}
	jsonBytes, _ := json.Marshal(updateData)
	req3, _ := http.NewRequest("PUT", server.URL+"/api/v1/media/"+att.ID, bytes.NewReader(jsonBytes))
	resp3, _ := http.DefaultClient.Do(req3)
	if resp3.StatusCode != http.StatusOK {
		t.Fatalf("PUT failed: %d", resp3.StatusCode)
	}

	var updated Attachment
	json.NewDecoder(resp3.Body).Decode(&updated)
	if updated.AltText == nil || *updated.AltText != "Beautiful landscape" {
		t.Errorf("alt text was not updated correctly: %v", updated.AltText)
	}
}
