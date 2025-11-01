package values

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func TestGetYamlPath(t *testing.T) {
	testDir := "testdata"

	type args struct {
		dirPath  string
		baseName string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name:    "Should return path for existing .yaml file",
			args:    args{dirPath: testDir, baseName: "file4"}, // file4.yaml exists in testdata
			want:    filepath.Join(testDir, "file4.yaml"),
			wantErr: nil,
		},
		{
			name:    "Should return path for existing .yml file",
			args:    args{dirPath: testDir, baseName: "file1"}, // file1.yml exists in testdata
			want:    filepath.Join(testDir, "file1.yml"),
			wantErr: nil,
		},
		{
			name:    "Should return error for non-existing file",
			args:    args{dirPath: testDir, baseName: "nonexistent"},
			want:    "",
			wantErr: errors.New("neither .yaml nor .yml file found for nonexistent in testdata"),
		},
		{
			name:    "Should return .yaml path when both .yaml and .yml exist (prefers .yaml)",
			args:    args{dirPath: testDir, baseName: "both"}, // both.yaml and both.yml exist in testdata
			want:    filepath.Join(testDir, "both.yaml"),      // Should return the .yaml file since it's checked first
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetYamlPath(tt.args.dirPath, tt.args.baseName)

			assertutils.AssertEqualErrors(t, tt.wantErr, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
