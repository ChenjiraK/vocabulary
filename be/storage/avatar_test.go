package storage

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func uploadHeader(t *testing.T, data []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("profile_picture", "../../untrusted.svg")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write(data); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest("POST", "/", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	if err := request.ParseMultipartForm(1 << 20); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = request.MultipartForm.RemoveAll() })
	return request.MultipartForm.File["profile_picture"][0]
}

func TestLocalAvatarStoreValidatesContentAndGeneratesNames(t *testing.T) {
	for _, format := range []string{"png", "jpeg"} {
		t.Run(format, func(t *testing.T) {
			var data bytes.Buffer
			picture := image.NewRGBA(image.Rect(0, 0, 2, 2))
			var err error
			if format == "png" {
				err = png.Encode(&data, picture)
			} else {
				err = jpeg.Encode(&data, picture, nil)
			}
			if err != nil {
				t.Fatal(err)
			}
			data.WriteString("untrusted trailing content")
			store := &LocalAvatarStore{Directory: t.TempDir()}
			header := uploadHeader(t, data.Bytes())
			first, err := store.Save(header)
			if err != nil {
				t.Fatal(err)
			}
			second, err := store.Save(header)
			if err != nil || first == second {
				t.Fatal("uploads must have unique names")
			}
			stored, err := os.ReadFile(filepath.Join(store.Directory, filepath.Base(first)))
			if err != nil || bytes.Contains(stored, []byte("untrusted")) {
				t.Fatal("stored image contains untrusted filename or content")
			}
			if _, actualFormat, err := image.Decode(bytes.NewReader(stored)); err != nil || actualFormat != format {
				t.Fatal("saved image has incorrect format")
			}
			if err := store.Delete(first); err != nil {
				t.Fatal(err)
			}
			if err := store.Delete(first); err != nil {
				t.Fatal("deleting an absent avatar should succeed")
			}
		})
	}
}

func TestLocalAvatarStoreRejectsInvalidUploads(t *testing.T) {
	var largeDimensions bytes.Buffer
	if err := png.Encode(&largeDimensions, image.NewRGBA(image.Rect(0, 0, 4097, 1))); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name string
		data []byte
		want error
	}{
		{"empty", nil, ErrInvalidAvatar},
		{"fake image", []byte("<svg>not a PNG</svg>"), ErrInvalidAvatar},
		{"truncated PNG", []byte("\x89PNG\r\n\x1a\n"), ErrInvalidAvatar},
		{"oversized dimensions", largeDimensions.Bytes(), ErrInvalidAvatar},
		{"oversized file", bytes.Repeat([]byte{'x'}, MaxAvatarBytes+1), ErrAvatarTooLarge},
	} {
		t.Run(tc.name, func(t *testing.T) {
			store := &LocalAvatarStore{Directory: t.TempDir()}
			header := uploadHeader(t, tc.data)
			// Verify the actual byte limit independently from the declared size.
			header.Size = 0
			if _, err := store.Save(header); !errors.Is(err, tc.want) {
				t.Fatalf("got error %v, want %v", err, tc.want)
			}
			files, err := os.ReadDir(store.Directory)
			if err != nil || len(files) != 0 {
				t.Fatal("invalid upload wrote files")
			}
		})
	}
}

func TestLocalAvatarStoreRejectsUnsafeDeletePaths(t *testing.T) {
	store := &LocalAvatarStore{Directory: t.TempDir()}
	for _, path := range []string{"", "/other/file.png", AvatarPrefix + "../file.png", AvatarPrefix + "..", AvatarPrefix + "nested/file.png", AvatarPrefix + "nested\\file.png"} {
		if err := store.Delete(path); err == nil {
			t.Fatalf("accepted unsafe path %q", path)
		}
	}
}

func TestPublicAvatarURL(t *testing.T) {
	t.Setenv("PUBLIC_ASSET_BASE_URL", "https://cdn.example.com/")
	if got := PublicAvatarURL(AvatarPrefix + "test.png"); got != "https://cdn.example.com/uploads/avatars/test.png" {
		t.Fatalf("unexpected public URL: %s", got)
	}
	if got := PublicAvatarURL("https://social.example.com/avatar.jpg"); got != "https://social.example.com/avatar.jpg" {
		t.Fatal("social avatar URL was changed")
	}
	t.Setenv("PUBLIC_ASSET_BASE_URL", "")
	t.Setenv("APP_PORT", "8000")
	if got := PublicAvatarURL(AvatarPrefix + "test.png"); !strings.HasPrefix(got, "http://localhost:8000/") {
		t.Fatalf("unexpected default URL: %s", got)
	}
}
