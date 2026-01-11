package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

type ImageData struct {
	Data        []byte
	MimeType    string
	Extension   string
	ContentType string
}

var (
	ErrInvalidDataURIFormat = errors.New("invalid data URI format: must be 'data:<mime>;base64,<data>'")
	ErrUnsupportedMimeType  = errors.New("unsupported MIME type: only image/png, image/jpeg, image/heic, image/heif are allowed")
	ErrBase64DecodeFailed   = errors.New("failed to decode base64 data")
)

var allowedMimeTypes = map[string]struct {
	Extension   string
	ContentType string
}{
	"image/png":  {Extension: ".png", ContentType: "image/png"},
	"image/jpeg": {Extension: ".jpg", ContentType: "image/jpeg"},
	"image/heic": {Extension: ".heic", ContentType: "image/heic"},
	"image/heif": {Extension: ".heic", ContentType: "image/heif"},
}

// DecodeImageDataURI parses a data URL, decodes the base64 payload,
// and returns image binary data with detected MIME type and extension.
// Supports PNG, JPEG, and HEIC formats.
func DecodeImageDataURI(dataURI string) (*ImageData, error) {
	mimeType, base64Data, err := parseDataURI(dataURI)
	if err != nil {
		return nil, err
	}

	typeInfo, ok := allowedMimeTypes[mimeType]
	if !ok {
		return nil, ErrUnsupportedMimeType
	}

	imageBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, ErrBase64DecodeFailed
	}

	finalMimeType := mimeType
	finalExtension := typeInfo.Extension
	finalContentType := typeInfo.ContentType

	// HEIC/HEIFの場合、ftypボックスからより正確な判定を試みる
	if mimeType == "image/heic" || mimeType == "image/heif" {
		if detectedMime := detectHEICFromFtyp(imageBytes); detectedMime != "" {
			finalMimeType = detectedMime
			if info, ok := allowedMimeTypes[detectedMime]; ok {
				finalExtension = info.Extension
				finalContentType = info.ContentType
			}
		}
	}

	return &ImageData{
		Data:        imageBytes,
		MimeType:    finalMimeType,
		Extension:   finalExtension,
		ContentType: finalContentType,
	}, nil
}

// parseDataURI extracts MIME type and base64 payload from a data URL.
// Expected format: data:<mime>;base64,<data>
func parseDataURI(dataURI string) (mimeType string, base64Data string, err error) {
	if !strings.HasPrefix(dataURI, "data:") {
		return "", "", ErrInvalidDataURIFormat
	}

	rest := strings.TrimPrefix(dataURI, "data:")

	commaIdx := strings.Index(rest, ",")
	if commaIdx == -1 {
		return "", "", ErrInvalidDataURIFormat
	}

	metaPart := rest[:commaIdx]
	dataPart := rest[commaIdx+1:]

	if !strings.Contains(metaPart, ";base64") {
		return "", "", ErrInvalidDataURIFormat
	}

	mimeType = strings.Split(metaPart, ";")[0]
	if mimeType == "" {
		return "", "", ErrInvalidDataURIFormat
	}

	return mimeType, dataPart, nil
}

// detectHEICFromFtyp checks ISO Base Media File Format (ISOBMFF) ftyp box
// to detect HEIC/HEIF format more accurately.
// Returns detected MIME type or empty string if detection fails.
func detectHEICFromFtyp(data []byte) string {
	// ftyp box structure:
	// - 4 bytes: box size
	// - 4 bytes: box type "ftyp"
	// - 4 bytes: major brand
	// - 4 bytes: minor version
	// - N*4 bytes: compatible brands
	if len(data) < 12 {
		return ""
	}

	// Check for "ftyp" at offset 4
	if string(data[4:8]) != "ftyp" {
		return ""
	}

	majorBrand := string(data[8:12])

	// HEIC brands
	heicBrands := map[string]bool{
		"heic": true,
		"heix": true,
		"hevc": true,
		"hevx": true,
		"heim": true,
		"heis": true,
		"hevm": true,
		"hevs": true,
		"mif1": true,
	}

	// HEIF brands (non-HEVC based)
	heifBrands := map[string]bool{
		"avif": true,
		"avis": true,
		"msf1": true,
	}

	if heicBrands[majorBrand] {
		return "image/heic"
	}
	if heifBrands[majorBrand] {
		return "image/heif"
	}

	return ""
}
