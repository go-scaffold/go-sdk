package templateproviders

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-scaffold/go-sdk/v2/pkg/filters"
	"github.com/stretchr/testify/assert"

	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/pasdam/go-utils/pkg/filetestutils"
)

func TestNewFileSystemProvider(t *testing.T) {
	type args struct {
		inputDir string
		filter   filters.Filter
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "should fill all the fields correctly",
			args: args{
				inputDir: "some-input-dir",
				filter:   filters.NewNoOpFilter(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFileSystemProvider(tt.args.inputDir, tt.args.filter).(*fileSystemProvider)

			assert.Equal(t, tt.args.inputDir, got.indexer.Dir)
			assert.Equal(t, tt.args.filter, got.filter)
		})
	}
}

func Test_fileSystemProvider_NextTemplate(t *testing.T) {
	type fields struct {
		dir     string
		fileter filters.Filter
	}
	type mocks struct {
		file      string
		filtered  bool
		openError error
	}
	type want struct {
		content string
		path    string
		err     error
	}
	tests := []struct {
		name   string
		fields fields
		mocks  mocks
		want   []want
	}{
		{
			name: "Should return EOF if dir has no files",
			fields: fields{
				dir: filetestutils.TempDir(t),
			},
			mocks: mocks{
				file:      "",
				filtered:  false,
				openError: nil,
			},
			want: []want{
				{
					err: io.EOF,
				},
			},
		},
		{
			name: "Should propagate error if one occur while indexing folder",
			fields: fields{
				dir: "",
			},
			mocks: mocks{
				file:      "",
				filtered:  false,
				openError: nil,
			},
			want: []want{
				{
					err: errors.New("open : no such file or directory"),
				},
			},
		},
		{
			name: "Should propagate error if one is thrown while opening the file",
			fields: fields{
				dir: filepath.Join("testdata", "file_system_provider"),
			},
			mocks: mocks{
				file:      "",
				filtered:  false,
				openError: errors.New("some-open-error"),
			},
			want: []want{
				{
					err: errors.New("some-open-error"),
				},
			},
		},
		{
			name: "Should process all files in the folder, when folder is relative to the current one",
			fields: fields{
				dir: filepath.Join("testdata", "file_system_provider"),
			},
			mocks: mocks{
				file:      "",
				filtered:  false,
				openError: nil,
			},
			want: []want{
				{
					content: "file0-content\n",
					path:    "file0",
				},
				{
					content: "file1-content\n",
					path:    "file1",
				},
				{
					content: "fileA-content\n",
					path:    filepath.Join("test_folder", "fileA"),
				},
				{
					err: io.EOF,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewFileSystemProvider(tt.fields.dir, &mockFilter{
				File:    tt.mocks.file,
				Include: !tt.mocks.filtered,
			})

			if tt.mocks.openError != nil {
				mockFile, err := os.Open(filepath.Join("testdata", "file_system_provider", "file0"))
				assert.NoError(t, err)
				mockOpen(t, "", mockFile, tt.mocks.openError)
			}

			for _, want := range tt.want {
				got, err := p.NextTemplate()

				assertutils.AssertEqualErrors(t, want.err, err)
				if want.err == nil {
					assert.NotNil(t, got)
					gotContent, err := ioutil.ReadAll(got.Reader)
					assert.NoError(t, err)
					assert.Equal(t, want.content, string(gotContent))
					assert.Equal(t, want.path, got.Path)
				} else {
					assert.Nil(t, got)
				}
			}
		})
	}
}

type mockFilter struct {
	File    string
	Include bool
}

func (m *mockFilter) Accept(filePath string) bool {
	if m.Include {
		return strings.Contains(filePath, m.File)

	} else {
		return !strings.Contains(filePath, m.File)
	}
}

func mockOpen(t *testing.T, expectedName string, file *os.File, err error) {
	originalValue := open
	open = func(name string) (*os.File, error) {
		// assert.Equal(t, expectedName, name)
		return file, err
	}
	t.Cleanup(func() { open = originalValue })
}
