package values

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYamlManifest(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr string
	}{
		{
			name: "Should load existing files with expected prefix",
			args: args{filepath.Join("testdata", "file1.yml")},
			want: map[string]interface{}{
				defaultManifestPrefix: map[string]interface{}{
					"key1": "value1",
				},
			},
			wantErr: "",
		},
		{
			name:    "Should return error if a file doesn't exist",
			args:    args{"non_existent_file.yml"},
			want:    nil,
			wantErr: "open non_existent_file.yml: no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadYamlManifest(tt.args.path)
			if len(tt.wantErr) > 0 {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadYamlManifest() = %v, want %v", got, tt.want)
			}
		})
	}
}
