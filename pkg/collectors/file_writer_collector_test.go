package collectors

import (
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/pasdam/go-utils/pkg/filetestutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewFileWriterCollector(t *testing.T) {
	type args struct {
		outDir        string
		nextCollector pipeline.Collector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should create collector without next one",
			args: args{
				outDir: "some-out-dir-without-collector",
			},
		},
		{
			name: "Should create collector with next one",
			args: args{
				outDir:        "some-out-dir-with-collector",
				nextCollector: &fileWriterCollector{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFileWriterCollector(tt.args.outDir, tt.args.nextCollector).(*fileWriterCollector)

			assert.Equal(t, tt.args.outDir, got.outDir)
			assert.Equal(t, tt.args.nextCollector, got.next)
		})
	}
}

func Test_fileWriterCollector_Collect(t *testing.T) {
	type fields struct {
		outDir string
	}
	type args struct {
		path    string
		content string
	}
	type want struct {
		err  error
		file bool
		next bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "should return no error if file is wrote correctly and no next collector is specified",
			fields: fields{
				outDir: filetestutils.TempDir(t),
			},
			args: args{
				path:    "some-correct-path",
				content: "some-correct-content",
			},
			want: want{
				err:  nil,
				file: true,
				next: false,
			},
		},
		{
			name: "should return no error if file is wrote correctly and next collector is invoked correctly",
			fields: fields{
				outDir: filetestutils.TempDir(t),
			},
			args: args{
				path:    "some-correct-path",
				content: "some-correct-content",
			},
			want: want{
				err:  nil,
				file: true,
				next: true,
			},
		},
		{
			name: "should propagate error if file is wrote correctly but next collector raises one",
			fields: fields{
				outDir: filetestutils.TempDir(t),
			},
			args: args{
				path:    "some-correct-path",
				content: "some-correct-content",
			},
			want: want{
				err:  nil,
				file: true,
				next: true,
			},
		},
		{
			name: "should return error if one is thrown while saving the output file",
			fields: fields{
				outDir: filepath.Join("testdata", "out", ".gitignore"),
			},
			args: args{
				path:    "some-not-existing-path",
				content: "some-content",
			},
			want: want{
				err:  errors.New("mkdir testdata/out/.gitignore: not a directory"),
				file: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var next pipeline.Collector
			if tt.want.next {
				next = &mockCollector{}
				next.(*mockCollector).On("Collect", mock.Anything).Return(tt.want.err)
			}
			p := &fileWriterCollector{
				baseCollector: baseCollector{
					next: next,
				},
				outDir: tt.fields.outDir,
			}

			err := p.Collect(&pipeline.Template{
				Path:   tt.args.path,
				Reader: io.NopCloser(strings.NewReader(tt.args.content)),
			})

			assertutils.AssertEqualErrors(t, tt.want.err, err)
			outPath := filepath.Join(tt.fields.outDir, tt.args.path)
			if tt.want.file {
				filetestutils.FileExistsWithContent(t, outPath, tt.args.content)
			} else {
				filetestutils.PathDoesNotExist(t, outPath)
			}
			if tt.want.next {
				next.(*mockCollector).AssertCalled(t, "Collect", mock.Anything)
			}
		})
	}
}
