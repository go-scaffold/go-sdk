package values

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func TestLoadYamlManifestAndValues(t *testing.T) {
	type args struct {
		manifestPath string
		valuesPaths  []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr error
	}{
		{
			name: "Should return merged manifest and values if no error is raised while loading the data",
			args: args{
				manifestPath: filepath.Join("testdata", "file1.yml"),
				valuesPaths: []string{
					filepath.Join("testdata", "file2.yml"),
					filepath.Join("testdata", "file3.yml"),
				},
			},
			want: map[string]interface{}{
				defaultManifestPrefix: map[string]interface{}{
					"key1": "value1",
				},
				defaultValuesPrefix: map[string]interface{}{
					"key2": "value2",
					"key1": "value3",
				},
			},
		},
		{
			name: "Should propagate error if load manifest throws one",
			args: args{
				manifestPath: "not-existing-manifest-file.yml",
				valuesPaths: []string{
					filepath.Join("testdata", "file2.yml"),
					filepath.Join("testdata", "file3.yml"),
				},
			},
			wantErr: errors.New("open not-existing-manifest-file.yml: no such file or directory"),
		},
		{
			name: "Should propagate error if load values throws one",
			args: args{
				manifestPath: filepath.Join("testdata", "file1.yml"),
				valuesPaths: []string{
					filepath.Join("testdata", "file2.yml"),
					"not-existing-values-file.yml",
					filepath.Join("testdata", "file3.yml"),
				},
			},
			wantErr: errors.New("open not-existing-values-file.yml: no such file or directory"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadYamlManifestAndValues(tt.args.manifestPath, tt.args.valuesPaths...)

			assert.Equal(t, tt.want, got)
			assertutils.AssertEqualErrors(t, tt.wantErr, err)
		})
	}
}
