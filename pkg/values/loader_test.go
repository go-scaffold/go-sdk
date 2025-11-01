package values

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pasdam/go-utils/pkg/filetestutils"
	"github.com/stretchr/testify/assert"
)

func TestNewLoader(t *testing.T) {
	loader := NewLoader()

	assert.Equal(t, "Manifest", loader.manifestPrefix)
	assert.Equal(t, "Manifest", loader.manifestBasename)
	assert.Equal(t, "Values", loader.valuesPrefix)
	assert.Equal(t, "values", loader.valuesBasename)
}

func TestNewLoaderWithValues(t *testing.T) {
	manifestPrefix := "CustomManifest"
	manifestBasename := "custom-manifest"
	valuesPrefix := "CustomValues"
	valuesBasename := "custom-values"

	loader := NewLoaderWithValues(manifestPrefix, manifestBasename, valuesPrefix, valuesBasename)

	assert.Equal(t, manifestPrefix, loader.manifestPrefix)
	assert.Equal(t, manifestBasename, loader.manifestBasename)
	assert.Equal(t, valuesPrefix, loader.valuesPrefix)
	assert.Equal(t, valuesBasename, loader.valuesBasename)
}

func TestLoader_LoadYAMLs(t *testing.T) {
	type args struct {
		createDir               bool
		manifestContent         string
		valuesContent           string
		additionalValueContents []string
		additionalValuePaths    []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr string
	}{
		{
			name: "Should load both manifest and values when both files exist",
			args: args{
				createDir:               true,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{},
			},
			want: map[string]interface{}{
				"Manifest": map[string]interface{}{
					"mk1": "mv1",
					"mk2": map[string]interface{}{
						"mk21": "mv21",
					},
				},
				"Values": map[string]interface{}{
					"vk1": "vv1",
					"vk2": map[string]interface{}{
						"vk21": "vv21",
					},
				},
			},
			wantErr: "",
		},
		{
			name: "Should load with additional value files when both manifest and values exist",
			args: args{
				createDir:               true,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{"ck1: cv1\nck2:\n  ck21: cv21\n"},
			},
			want: map[string]interface{}{
				"Manifest": map[string]interface{}{
					"mk1": "mv1",
					"mk2": map[string]interface{}{
						"mk21": "mv21",
					},
				},
				"Values": map[string]interface{}{
					"vk1": "vv1",
					"vk2": map[string]interface{}{
						"vk21": "vv21",
					},
					"ck1": "cv1",
					"ck2": map[string]interface{}{
						"ck21": "cv21",
					},
				},
			},
			wantErr: "",
		},
		{
			name: "Should return error when manifest file doesn't exist",
			args: args{
				createDir:               true,
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{},
			},
			want:    nil,
			wantErr: "an error occurred while getting the manifest path: neither .yaml nor .yml file found for Manifest in ",
		},
		{
			name: "Should return error when manifest is not a valid YAML",
			args: args{
				createDir:               true,
				manifestContent:         "this, is, not, valid",
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{},
			},
			want:    nil,
			wantErr: "an error occurred while reading the manifest file: failed to parse yaml: yaml: unmarshal errors",
		},
		{
			name: "Should return error when values file doesn't exist",
			// GetYamlValuesPath returns empty string when no values exists,
			// which causes LoadYamlFilesWithPrefix to fail as it tries to load empty path
			args: args{
				createDir:               true,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				additionalValueContents: []string{},
			},
			want:    nil,
			wantErr: "an error occurred while getting the value path: neither .yaml nor .yml file found for values in ",
		},
		{
			name: "Should return error when values file is not a valid YAML",
			// GetYamlValuesPath returns empty string when no values exists,
			// which causes LoadYamlFilesWithPrefix to fail as it tries to load empty path
			args: args{
				createDir:               true,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				valuesContent:           "this, is, not, valid",
				additionalValueContents: []string{},
			},
			want:    nil,
			wantErr: "error while loading data: failed to parse yaml: yaml: unmarshal errors",
		},
		{
			name: "Should return error when template root path doesn't exist",
			args: args{
				createDir:               false,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{},
			},
			want:    nil,
			wantErr: "an error occurred while getting the manifest path: neither .yaml nor .yml file found",
		},
		{
			name: "Should return error when additional value path doesn't exist",
			args: args{
				createDir:               true,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{},
				additionalValuePaths:    []string{"not_existing_path.yaml"},
			},
			want:    nil,
			wantErr: "error while loading data: open not_existing_path.yaml: no such file or directory",
		},
		{
			name: "Should return error when additional value file is not a valid YAML",
			args: args{
				createDir:               true,
				manifestContent:         "mk1: mv1\nmk2:\n  mk21: mv21\n",
				valuesContent:           "vk1: vv1\nvk2:\n  vk21: vv21\n",
				additionalValueContents: []string{"this, is, not, valid"},
			},
			want:    nil,
			wantErr: "error while loading data: failed to parse yaml: yaml: unmarshal errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLoader()
			var dir string
			valuePaths := make([]string, 0, len(tt.args.additionalValueContents)+len(tt.args.additionalValuePaths))
			valuePaths = append(valuePaths, tt.args.additionalValuePaths...)
			if tt.args.createDir {
				dir = filetestutils.TempDir(t)
				if len(tt.args.manifestContent) > 0 {
					err := os.WriteFile(filepath.Join(dir, l.manifestBasename+".yml"), []byte(tt.args.manifestContent), 0644)
					assert.NoError(t, err)
				}
				if len(tt.args.valuesContent) > 0 {
					err := os.WriteFile(filepath.Join(dir, l.valuesBasename+".yml"), []byte(tt.args.valuesContent), 0644)
					assert.NoError(t, err)
				}
				for i, content := range tt.args.additionalValueContents {
					path := filepath.Join(dir, fmt.Sprintf("%s-%d.yml", l.valuesBasename, i))
					valuePaths = append(valuePaths, path)
					err := os.WriteFile(path, []byte(content), 0644)
					assert.NoError(t, err)
				}
			} else {
				dir = "not_existing_dir"
			}

			got, err := l.LoadYAMLs(dir, valuePaths)
			if len(tt.wantErr) > 0 {
				assert.Contains(t, err.Error(), tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
