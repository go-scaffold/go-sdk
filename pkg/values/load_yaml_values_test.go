package values

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYamlValues(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr string
	}{
		{
			name: "Should load existing files with expected prefix",
			args: args{paths: []string{filepath.Join("testdata", "file1.yml"), filepath.Join("testdata", "file2.yml")}},
			want: map[string]interface{}{
				defaultValuesPrefix: map[string]interface{}{
					"key1": "value1",
					"key2": "value2",
				},
			},
			wantErr: "",
		},
		{
			name: "Should load existing files with expected prefix, and override values",
			args: args{paths: []string{filepath.Join("testdata", "file1.yml"), filepath.Join("testdata", "file3.yml")}},
			want: map[string]interface{}{
				defaultValuesPrefix: map[string]interface{}{
					"key1": "value3",
				},
			},
			wantErr: "",
		},
		{
			name:    "Should return error if a file doesn't exist",
			args:    args{paths: []string{"non_existent_file.yml"}},
			want:    nil,
			wantErr: "open non_existent_file.yml: no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadYamlValues(tt.args.paths)
			if len(tt.wantErr) > 0 {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			} else {
				assert.NoError(t, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadYamlValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
