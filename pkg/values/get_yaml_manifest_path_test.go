package values

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func TestGetYamlManifestPath(t *testing.T) {
	testCases := []struct {
		name    string
		args    struct{ dir string }
		want    string
		wantErr error
	}{
		{
			name:    "Should return path for existing Manifest.yaml file (takes priority over .yml)",
			args:    struct{ dir string }{dir: "testdata"},
			want:    filepath.Join("testdata", "Manifest.yaml"), // Should return .yaml since it has priority
			wantErr: nil,
		},
		{
			name:    "Should return error when directory doesn't exist",
			args:    struct{ dir string }{dir: "/nonexistent/directory/path"},
			want:    "",
			wantErr: errors.New("neither .yaml nor .yml file found for Manifest in /nonexistent/directory/path"),
		},
		{
			name:    "Should return error when neither Manifest.yaml nor Manifest.yml exist",
			args:    struct{ dir string }{dir: "."},
			want:    "",
			wantErr: errors.New("neither .yaml nor .yml file found for Manifest in ."),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetYamlManifestPath(tt.args.dir)

			assertutils.AssertEqualErrors(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
