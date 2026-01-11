package tests

import (
	"os"
	"testing"

	"github.com/hackathon-20260110/api/adapter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeObjectKey(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "simple path",
			input:       "images/photo.jpg",
			expected:    "images/photo.jpg",
			expectError: false,
		},
		{
			name:        "leading slash removed",
			input:       "/images/photo.jpg",
			expected:    "images/photo.jpg",
			expectError: false,
		},
		{
			name:        "multiple leading slashes removed",
			input:       "///images/photo.jpg",
			expected:    "images/photo.jpg",
			expectError: false,
		},
		{
			name:        "empty path rejected",
			input:       "",
			expected:    "",
			expectError: true,
			errorMsg:    "object key cannot be empty",
		},
		{
			name:        "only slashes rejected",
			input:       "///",
			expected:    "",
			expectError: true,
			errorMsg:    "object key cannot be empty",
		},
		{
			name:        "path traversal rejected",
			input:       "images/../secret.txt",
			expected:    "",
			expectError: true,
			errorMsg:    "object key cannot contain '..'",
		},
		{
			name:        "double dots at start rejected",
			input:       "../etc/passwd",
			expected:    "",
			expectError: true,
			errorMsg:    "object key cannot contain '..'",
		},
		{
			name:        "file with dots allowed",
			input:       "images/photo.2024.01.jpg",
			expected:    "images/photo.2024.01.jpg",
			expectError: false,
		},
		{
			name:        "single dot allowed",
			input:       "./images/photo.jpg",
			expected:    "./images/photo.jpg",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := adapter.NormalizeObjectKey(tt.input)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestBuildPublicURL(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		objectKey string
		expected  string
	}{
		{
			name:      "simple URL construction",
			baseURL:   "https://pub-xxx.r2.dev",
			objectKey: "images/photo.jpg",
			expected:  "https://pub-xxx.r2.dev/images/photo.jpg",
		},
		{
			name:      "custom domain",
			baseURL:   "https://cdn.example.com",
			objectKey: "uploads/2024/photo.png",
			expected:  "https://cdn.example.com/uploads/2024/photo.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("R2_PUBLIC_BASE_URL", tt.baseURL)
			defer os.Unsetenv("R2_PUBLIC_BASE_URL")

			result := adapter.BuildPublicURL(tt.objectKey)
			assert.Equal(t, tt.expected, result)
		})
	}
}
