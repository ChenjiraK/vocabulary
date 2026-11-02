package storage

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	MaxAvatarBytes = 2 << 20
	AvatarPrefix   = "/uploads/avatars/"
)

var (
	ErrAvatarTooLarge = errors.New("profile_picture must be at most 2 MB")
	ErrInvalidAvatar  = errors.New("profile_picture must be a valid JPEG or PNG image up to 4096 x 4096 pixels")
)

type AvatarStore interface {
	Save(file *multipart.FileHeader) (string, error)
	Delete(path string) error
}

type LocalAvatarStore struct {
	Directory string
}

func AvatarDirectory() string {
	if directory := os.Getenv("AVATAR_UPLOAD_DIR"); directory != "" {
		return directory
	}
	return "./uploads/avatars"
}

func PublicAvatarURL(path string) string {
	if u, err := url.Parse(path); err == nil && u.IsAbs() {
		return path
	}
	base := strings.TrimRight(os.Getenv("PUBLIC_ASSET_BASE_URL"), "/")
	if base == "" {
		port := os.Getenv("APP_PORT")
		if port == "" {
			port = "8080"
		}
		base = "http://localhost:" + port
	}
	return base + "/" + strings.TrimLeft(path, "/")
}

func (s *LocalAvatarStore) Save(file *multipart.FileHeader) (string, error) {
	if file.Size > MaxAvatarBytes {
		return "", ErrAvatarTooLarge
	}
	source, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("open avatar: %w", err)
	}
	defer source.Close()
	data, err := io.ReadAll(io.LimitReader(source, MaxAvatarBytes+1))
	if err != nil {
		return "", fmt.Errorf("read avatar: %w", err)
	}
	if len(data) > MaxAvatarBytes {
		return "", ErrAvatarTooLarge
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "png" && format != "jpeg") || config.Width < 1 || config.Height < 1 || config.Width > 4096 || config.Height > 4096 {
		return "", ErrInvalidAvatar
	}
	picture, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", ErrInvalidAvatar
	}
	var id [16]byte
	if _, err := rand.Read(id[:]); err != nil {
		return "", fmt.Errorf("generate avatar name: %w", err)
	}
	extension := ".png"
	if format == "jpeg" {
		extension = ".jpg"
	}
	name := hex.EncodeToString(id[:]) + extension
	if err := os.MkdirAll(s.Directory, 0755); err != nil {
		return "", fmt.Errorf("create avatar directory: %w", err)
	}
	destination := filepath.Join(s.Directory, name)
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return "", fmt.Errorf("create avatar: %w", err)
	}
	// Re-encode the image so uploaded metadata and trailing content are discarded.
	if format == "jpeg" {
		err = jpeg.Encode(output, picture, &jpeg.Options{Quality: 90})
	} else {
		err = png.Encode(output, picture)
	}
	closeErr := output.Close()
	if err != nil || closeErr != nil {
		_ = os.Remove(destination)
		return "", fmt.Errorf("write avatar: %w", errors.Join(err, closeErr))
	}
	return AvatarPrefix + name, nil
}

func (s *LocalAvatarStore) Delete(path string) error {
	name, ok := strings.CutPrefix(path, AvatarPrefix)
	if !ok || name == "" || name == "." || name == ".." || strings.ContainsAny(name, "/\\") {
		return errors.New("invalid avatar path")
	}
	if err := os.Remove(filepath.Join(s.Directory, name)); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	return nil
}
