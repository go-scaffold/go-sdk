package values

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func TestGetYamlValuesPath(t *testing.T) {
	testCases := []struct {
		name    string
		args    struct{ dir string }
		want    string
		wantErr error
	}{
		{
			name:    "Should return path for existing Values.yaml file (takes priority over .yml)",
			args:    struct{ dir string }{dir: "testdata"},
			want:    filepath.Join("testdata", "Values.yaml"), // Should return .yaml since it has priority
			wantErr: nil,
		},
		{
			name:    "Should return empty string when directory doesn't exist (not an error case)",
			args:    struct{ dir string }{dir: "/nonexistent/directory/path"},
			want:    "",
			wantErr: errors.New("neither .yaml nor .yml file found for Values in /nonexistent/directory/path"),
		},
		{
			name:    "Should return empty string when neither Values.yaml nor Values.yml exist",
			args:    struct{ dir string }{dir: "."},
			want:    "",
			wantErr: errors.New("neither .yaml nor .yml file found for Values in ."),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetYamlValuesPath(tt.args.dir)

			assertutils.AssertEqualErrors(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
