package collectors

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-io-utilx/pkg/ioutilx"
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

			assert.Equal(t, tt.args.outDir, got.opts.OutDir)
			assert.Equal(t, tt.args.nextCollector, got.next)
			assert.Equal(t, false, got.opts.SkipUnchanged)
		})
	}
}

func TestNewFileWriterCollectorWithOpts(t *testing.T) {
	type args struct {
		opts          FileWriterCollectorOptions
		nextCollector pipeline.Collector
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Should create collector without next one",
			args: args{
				opts: FileWriterCollectorOptions{
					OutDir: "some-out-dir-without-collector",
				},
			},
		},
		{
			name: "Should create collector with next one",
			args: args{
				opts: FileWriterCollectorOptions{
					OutDir: "some-out-dir-with-collector",
				},
				nextCollector: &fileWriterCollector{},
			},
		},
		{
			name: "Should create collector to skip unchanged files",
			args: args{
				opts: FileWriterCollectorOptions{
					OutDir:        "some-out-dir-skip-unchanged",
					SkipUnchanged: true,
				},
				nextCollector: &fileWriterCollector{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFileWriterCollectorWithOpts(tt.args.opts, tt.args.nextCollector).(*fileWriterCollector)

			assert.Equal(t, tt.args.opts.OutDir, got.opts.OutDir)
			assert.Equal(t, tt.args.nextCollector, got.next)
			assert.Equal(t, tt.args.opts.SkipUnchanged, got.opts.SkipUnchanged)
		})
	}
}

func Test_fileWriterCollector_Collect(t *testing.T) {
	type args struct {
		path    string
		content string
	}
	type want struct {
		err       error
		file      bool
		next      bool
		overwrite bool
	}
	tests := []struct {
		name     string
		opts     FileWriterCollectorOptions
		mockFile bool
		args     args
		want     want
	}{
		{
			name: "should return no error if file is wrote correctly and no next collector is specified",
			opts: FileWriterCollectorOptions{
				OutDir: filetestutils.TempDir(t),
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
			name: "should overwrite file if SkipUnchanged is false",
			opts: FileWriterCollectorOptions{
				OutDir: filetestutils.TempDir(t),
			},
			args: args{
				path:    "some-correct-path",
				content: "some-correct-content",
			},
			mockFile: true,
			want: want{
				err:       nil,
				file:      true,
				next:      false,
				overwrite: true,
			},
		},
		{
			name: "should overwrite file if SkipUnchanged is true",
			opts: FileWriterCollectorOptions{
				OutDir:        filetestutils.TempDir(t),
				SkipUnchanged: true,
			},
			args: args{
				path:    "some-correct-path",
				content: "some-correct-content",
			},
			mockFile: true,
			want: want{
				err:       nil,
				file:      true,
				next:      false,
				overwrite: false,
			},
		},
		{
			name: "should return no error if file is wrote correctly and next collector is invoked correctly",
			opts: FileWriterCollectorOptions{
				OutDir: filetestutils.TempDir(t),
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
			opts: FileWriterCollectorOptions{
				OutDir: filetestutils.TempDir(t),
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
			opts: FileWriterCollectorOptions{
				OutDir: filepath.Join("testdata", "out", ".gitignore"),
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
				next.(*mockCollector).On("OnPipelineCompleted").Return(nil) // Expect the new method to be called
			}
			p := NewFileWriterCollectorWithOpts(tt.opts, next)
			beforeTimestamp := int64(0)
			outPath := filepath.Join(tt.opts.OutDir, tt.args.path)
			if tt.mockFile {
				ioutilx.ReaderToFile(bytes.NewReader([]byte(tt.args.content)), outPath)
				stat, err := os.Stat(outPath)
				assert.NoError(t, err)
				beforeTimestamp = stat.ModTime().UnixNano()
			}

			err := p.Collect(&pipeline.Template{
				Path:   tt.args.path,
				Reader: io.NopCloser(strings.NewReader(tt.args.content)),
			})

			assertutils.AssertEqualErrors(t, tt.want.err, err)
			if tt.want.file {
				if tt.mockFile {
					stat, err := os.Stat(outPath)
					assert.NoError(t, err)
					currentTimestamp := stat.ModTime().UnixNano()

					if tt.want.overwrite {
						assert.Greater(t, currentTimestamp, beforeTimestamp)
					} else {
						assert.Equal(t, beforeTimestamp, currentTimestamp)
					}
				}

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

func Test_fileWriterCollector_OnPipelineCompleted(t *testing.T) {
	tests := []struct {
		name             string
		nextError        error
		expectedError    error
		nextExists       bool
		cleanupUntracked bool
		setupFiles       map[string]string // path -> content for files to create before test
		expectedFiles    []string          // files that should exist after completion
		unexpectedFiles  []string          // files that should NOT exist after completion
	}{
		{
			name:             "Should return nil when no next collector exists",
			nextExists:       false,
			expectedError:    nil,
			cleanupUntracked: false,
		},
		{
			name:             "Should return nil when next collector returns nil",
			nextExists:       true,
			nextError:        nil,
			expectedError:    nil,
			cleanupUntracked: false,
		},
		{
			name:             "Should return error when next collector returns error",
			nextExists:       true,
			nextError:        errors.New("next-collector-error"),
			expectedError:    errors.New("next-collector-error"),
			cleanupUntracked: false,
		},
		{
			name:             "Should not remove untracked files when cleanup is disabled",
			nextExists:       false,
			expectedError:    nil,
			cleanupUntracked: false,
			setupFiles: map[string]string{
				"existing-file.txt": "existing content",
				"tracked-file.txt":  "tracked content",
			},
			expectedFiles: []string{"existing-file.txt", "tracked-file.txt"},
		},
		{
			name:             "Should remove untracked files when cleanup is enabled",
			nextExists:       false,
			expectedError:    nil,
			cleanupUntracked: true,
			setupFiles: map[string]string{
				"generated-file.txt": "generated content", // This will be tracked in the test
				"untracked-file.txt": "untracked content", // This should be removed
			},
			expectedFiles:   []string{"generated-file.txt"}, // Only the generated file should remain
			unexpectedFiles: []string{"untracked-file.txt"}, // The untracked file should be removed
		},
		{
			name:             "Should preserve generated files and remove untracked when cleanup is enabled",
			nextExists:       false,
			expectedError:    nil,
			cleanupUntracked: true,
			setupFiles: map[string]string{
				"generated1.txt":  "content1",    // This will be tracked
				"generated2.txt":  "content2",    // This will be tracked
				"preexisting.txt": "old content", // This was there before and should be removed if not tracked
			},
			expectedFiles:   []string{"generated1.txt", "generated2.txt"}, // Only generated files should remain
			unexpectedFiles: []string{"preexisting.txt"},                  // Pre-existing file should be removed
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary output directory
			tempDir := filetestutils.TempDir(t)

			// Set up initial files if needed
			for filename, content := range tt.setupFiles {
				fullPath := filepath.Join(tempDir, filename)
				err := ioutilx.ReaderToFile(bytes.NewReader([]byte(content)), fullPath)
				assert.NoError(t, err)
			}

			var nextCollector pipeline.Collector
			if tt.nextExists {
				nextCollector = &mockCollector{}
				nextCollector.(*mockCollector).On("OnPipelineCompleted").Return(tt.nextError)
			}

			// Create collector with options that include cleanup setting
			opts := FileWriterCollectorOptions{
				OutDir:           tempDir,
				CleanupUntracked: tt.cleanupUntracked,
			}
			p := NewFileWriterCollectorWithOpts(opts, nextCollector).(*fileWriterCollector)

			// Simulate that during pipeline execution, files were generated
			// (added to generatedFiles map) - for the cleanup test cases
			if tt.cleanupUntracked {
				for _, expectedFile := range tt.expectedFiles {
					fullPath := filepath.Join(tempDir, expectedFile)
					p.generatedFiles[fullPath] = true
				}
			}

			err := p.OnPipelineCompleted()

			assertutils.AssertEqualErrors(t, tt.expectedError, err)
			if tt.nextExists {
				nextCollector.(*mockCollector).AssertCalled(t, "OnPipelineCompleted")
			}

			// Verify expected files still exist
			for _, expectedFile := range tt.expectedFiles {
				fullPath := filepath.Join(tempDir, expectedFile)
				_, err := os.Stat(fullPath)
				assert.False(t, os.IsNotExist(err), "File %s should exist but it doesn't", fullPath)
			}

			// Verify unexpected files were removed
			for _, unexpectedFile := range tt.unexpectedFiles {
				fullPath := filepath.Join(tempDir, unexpectedFile)
				_, err := os.Stat(fullPath)
				assert.True(t, os.IsNotExist(err), "File %s should not exist but it does", fullPath)
			}
		})
	}
}
