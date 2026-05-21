package media

import (
	"bytes"
	"context"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"

	"github.com/buckket/go-blurhash"
	"github.com/sumi-devs/canopy-social/canopy/pkg/config"
	"github.com/sumi-devs/canopy-social/canopy/pkg/ulid"
)

type S3Client interface {
	Upload(ctx context.Context, key string, contentType string, reader io.Reader) (string, error)
	Delete(ctx context.Context, key string) error
}

type Service struct {
	repo Repository
	s3   S3Client
	cfg  *config.Config
}

func NewService(repo Repository, s3 S3Client, cfg *config.Config) *Service {
	return &Service{
		repo: repo,
		s3:   s3,
		cfg:  cfg,
	}
}

func (s *Service) ProcessUpload(ctx context.Context, accountID string, originalReader io.Reader, filename string) (*Attachment, error) {
	var fileBuf bytes.Buffer
	teeReader := io.TeeReader(originalReader, &fileBuf)

	header := make([]byte, 512)
	n, _ := io.ReadFull(teeReader, header)
	if n == 0 {
		return nil, fmt.Errorf("empty file uploaded")
	}

	mimeType := http.DetectContentType(header)

	_, err := io.Copy(&fileBuf, originalReader)
	if err != nil {
		return nil, fmt.Errorf("failed to read full file data: %w", err)
	}

	fileSize := int64(fileBuf.Len())

	mediaType := "unknown"
	if strings.HasPrefix(mimeType, "image/") {
		mediaType = "image"
	} else if strings.HasPrefix(mimeType, "video/") {
		mediaType = "video"
	} else if strings.HasPrefix(mimeType, "audio/") {
		mediaType = "audio"
	}

	if mediaType == "image" && s.cfg.Features.MaxImageMB > 0 {
		if fileSize > int64(s.cfg.Features.MaxImageMB)*1024*1024 {
			return nil, fmt.Errorf("image file size exceeds maximum allowed size of %d MB", s.cfg.Features.MaxImageMB)
		}
	} else if mediaType == "video" && s.cfg.Features.MaxVideoMB > 0 {
		if fileSize > int64(s.cfg.Features.MaxVideoMB)*1024*1024 {
			return nil, fmt.Errorf("video file size exceeds maximum allowed size of %d MB", s.cfg.Features.MaxVideoMB)
		}
	}

	attachmentID := ulid.New()
	att := &Attachment{
		ID:          attachmentID,
		AccountID:   accountID,
		Type:        mediaType,
		MimeType:    &mimeType,
		FileSize:    &fileSize,
		IsProcessed: false,
		CreatedAt:   time.Now(),
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if ext == "" {
		ext = ".bin"
		if mediaType == "image" {
			ext = ".jpg"
		}
	}

	originalKey := fmt.Sprintf("accounts/%s/media/%s%s", accountID, attachmentID, ext)
	att.StorageKey = &originalKey

	originalURL, err := s.s3.Upload(ctx, originalKey, mimeType, bytes.NewReader(fileBuf.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("failed to upload original file to S3: %w", err)
	}
	att.URL = &originalURL

	if mediaType == "image" {
		img, _, err := image.Decode(bytes.NewReader(fileBuf.Bytes()))
		if err == nil {
			bounds := img.Bounds()
			w := bounds.Dx()
			h := bounds.Dy()
			att.Width = &w
			att.Height = &h

			if hashStr, err := blurhash.Encode(4, 3, img); err == nil {
				att.Blurhash = &hashStr
			}

			thumbImg := resizeImage(img, 400)
			var thumbBuf bytes.Buffer
			if err := jpeg.Encode(&thumbBuf, thumbImg, nil); err == nil {
				thumbKey := fmt.Sprintf("accounts/%s/media/%s_thumb.jpg", accountID, attachmentID)
				thumbURL, err := s.s3.Upload(ctx, thumbKey, "image/jpeg", bytes.NewReader(thumbBuf.Bytes()))
				if err == nil {
					att.ThumbnailURL = &thumbURL
				}
			}
		}
	}

	att.IsProcessed = true
	saved, err := s.repo.Create(ctx, att)
	if err != nil {
		return nil, fmt.Errorf("failed to save media metadata: %w", err)
	}

	return saved, nil
}

func (s *Service) GetAttachment(ctx context.Context, id string) (*Attachment, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) UpdateAttachment(ctx context.Context, id string, altText *string) (*Attachment, error) {
	return s.repo.Update(ctx, id, altText)
}

func (s *Service) DeleteAttachment(ctx context.Context, id string) error {
	att, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if att.StorageKey != nil {
		s.s3.Delete(ctx, *att.StorageKey)
	}

	if att.ThumbnailURL != nil {
		parts := strings.Split(*att.ThumbnailURL, "/")
		if len(parts) > 0 {
			thumbKey := strings.Join(parts[len(parts)-4:], "/")
			if strings.Contains(thumbKey, "_thumb.jpg") {
				s.s3.Delete(ctx, thumbKey)
			}
		}
	}

	return s.repo.Delete(ctx, id)
}

func resizeImage(img image.Image, maxDim int) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		return img
	}

	scale := float64(maxDim) / math.Max(float64(w), float64(h))
	newW := int(float64(w) * scale)
	newH := int(float64(h) * scale)

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.BiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}
