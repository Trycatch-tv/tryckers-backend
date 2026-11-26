package storage

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type mockS3Client struct {
	objects   map[string][]byte
	types     map[string]string
	putErr    error
	getErr    error
	deleteErr error
}

func newMockS3Client() *mockS3Client {
	return &mockS3Client{
		objects: make(map[string][]byte),
		types:   make(map[string]string),
	}
}

func (m *mockS3Client) PutObject(_ context.Context, params *s3.PutObjectInput, _ ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if m.putErr != nil {
		return nil, m.putErr
	}
	data, err := io.ReadAll(params.Body)
	if err != nil {
		return nil, err
	}
	m.objects[aws.ToString(params.Key)] = data
	m.types[aws.ToString(params.Key)] = aws.ToString(params.ContentType)
	return &s3.PutObjectOutput{}, nil
}

func (m *mockS3Client) GetObject(_ context.Context, params *s3.GetObjectInput, _ ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	key := aws.ToString(params.Key)
	data, ok := m.objects[key]
	if !ok {
		return nil, errors.New("NoSuchKey")
	}
	ct := m.types[key]
	return &s3.GetObjectOutput{
		Body:        io.NopCloser(bytes.NewReader(data)),
		ContentType: aws.String(ct),
	}, nil
}

func (m *mockS3Client) DeleteObject(_ context.Context, params *s3.DeleteObjectInput, _ ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if m.deleteErr != nil {
		return nil, m.deleteErr
	}
	delete(m.objects, aws.ToString(params.Key))
	return &s3.DeleteObjectOutput{}, nil
}

func createMultipartFileHeader(t *testing.T, filename string, payload []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write(payload); err != nil {
		t.Fatalf("failed to write payload: %v", err)
	}
	_ = writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err := req.ParseMultipartForm(10 << 20); err != nil {
		t.Fatalf("failed to parse multipart form: %v", err)
	}
	_, header, err := req.FormFile("file")
	if err != nil {
		t.Fatalf("failed to get form file: %v", err)
	}
	return header
}

// Valid 1x1 PNG bytes
var samplePNG = []byte{
	0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
	0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
	0x89, 0x00, 0x00, 0x00, 0x0A, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE,
	0x42, 0x60, 0x82,
}

func TestS3Storage_SaveGetDeleteAndResolve(t *testing.T) {
	mock := newMockS3Client()
	svc := &S3Storage{
		Client:        mock,
		Bucket:        "tryckers-media-storage-dev",
		PublicBaseURL: "https://api.tryckers.trycatch.tv",
	}

	fh := createMultipartFileHeader(t, "avatar.png", samplePNG)
	key, err := svc.SaveImage(context.Background(), fh, SaveOptions{
		Kind:     AvatarKind,
		OwnerID:  "user-123",
		MaxBytes: 2 * 1024 * 1024,
	})
	if err != nil {
		t.Fatalf("SaveImage failed: %v", err)
	}
	if !strings.HasPrefix(key, "avatars/user-123/") || !strings.HasSuffix(key, ".png") {
		t.Fatalf("unexpected object key format: %s", key)
	}

	resolved := svc.ResolveURL(key)
	if resolved != "https://api.tryckers.trycatch.tv/"+key {
		t.Fatalf("unexpected resolved URL: %s", resolved)
	}

	// Legacy URLs must NOT be resolved
	if svc.ResolveURL("/uploads/avatars/old.png") != "" || svc.ResolveURL("https://old.example.com/a.png") != "" {
		t.Fatalf("legacy URLs must return empty string")
	}

	reader, contentType, err := svc.GetObject(context.Background(), key)
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}
	defer reader.Close()
	if contentType != "image/png" {
		t.Fatalf("expected image/png, got %s", contentType)
	}

	if err := svc.DeleteObject(context.Background(), key); err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}
	if _, _, err := svc.GetObject(context.Background(), key); err == nil {
		t.Fatalf("expected error after deleting object")
	}
}

func TestS3Storage_RejectsInvalidTypeAndOversize(t *testing.T) {
	mock := newMockS3Client()
	svc := &S3Storage{Client: mock, Bucket: "test-bucket"}

	textHeader := createMultipartFileHeader(t, "bad.txt", []byte("not an image at all"))
	if _, err := svc.SaveImage(context.Background(), textHeader, SaveOptions{
		Kind:     PostKind,
		OwnerID:  "u1",
		MaxBytes: 1024,
	}); err == nil {
		t.Fatalf("expected error for invalid MIME type")
	}

	pngHeader := createMultipartFileHeader(t, "small.png", samplePNG)
	if _, err := svc.SaveImage(context.Background(), pngHeader, SaveOptions{
		Kind:     AvatarKind,
		OwnerID:  "u1",
		MaxBytes: 10,
	}); err == nil {
		t.Fatalf("expected error for oversized image")
	}
}
