package values

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetYamlManifestPath(t *testing.T) {
	// For the case where no manifest files exist, use a temp directory
	tempDir := t.TempDir()

	// Define test cases
	testCases := []struct {
		name    string
		args    struct{ dir string }
		want    string
		wantErr bool
	}{
		{
			name:    "Should return path for existing Manifest.yaml file (takes priority over .yml)",
			args:    struct{ dir string }{dir: "testdata"},
			want:    filepath.Join("testdata", "Manifest.yaml"), // Should return .yaml since it has priority
			wantErr: false,
		},
		{
			name:    "Should return empty string when directory doesn't exist (not an error case)",
			args:    struct{ dir string }{dir: "/nonexistent/directory/path"},
			want:    "",
			wantErr: false, // No error - just no file found
		},
		{
			name:    "Should return empty string when neither Manifest.yaml nor Manifest.yml exist",
			args:    struct{ dir string }{dir: tempDir},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetYamlManifestPath(tt.args.dir)
			if (err != nil) != tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err != nil)
			} else {
				assert.NoError(t, err)
			}
			if got != tt.want {
				t.Errorf("GetYamlManifestPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
