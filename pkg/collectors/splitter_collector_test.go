package collectors

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_splitterCollector_Collect(t *testing.T) {
	type mocks struct {
		nextCollectResult []error
	}
	type data struct {
		path    string
		content string
	}
	tests := []struct {
		name    string
		mocks   mocks
		args    data
		want    []data
		wantErr error
	}{
		{
			name: "Should accept result if file starts with right prefix and extract the only file",
			mocks: mocks{
				nextCollectResult: []error{
					nil,
				},
			},
			args: data{
				path:    "some-path/mul_something",
				content: "@@ name=\"some-other-name-1\"\nsome-file-1-line-1",
			},
			want: []data{
				{
					path:    "some-other-name-1",
					content: "some-file-1-line-1",
				},
			},
			wantErr: nil,
		},
		{
			name: "Should accept result if file starts with right prefix and extract separate files",
			mocks: mocks{
				nextCollectResult: []error{
					nil,
					nil,
				},
			},
			args: data{
				path:    "some-path/mul_something",
				content: "@@ name=\"some-other-name-1\"\nsome-file-1-line-1\nsome-file-1-line-2\n@@ name=\"some-other-name-2\"\nsome-file-2-line-1\n\n",
			},
			want: []data{
				{
					path:    "some-other-name-1",
					content: "some-file-1-line-1\nsome-file-1-line-2\n",
				},
				{
					path:    "some-other-name-2",
					content: "some-file-2-line-1\n\n",
				},
			},
			wantErr: nil,
		},
		{
			name: "Should not accept file if the name prefix is not the expected one",
			mocks: mocks{
				nextCollectResult: []error{
					nil,
				},
			},
			args: data{
				path:    "some-path/something",
				content: "@@ name=\"some-other-name-1\"\nsome-content-1\n@@ name=\"some-other-name-2\"\nsome-content-2",
			},
			want: []data{
				{
					path:    "some-path/something",
					content: "@@ name=\"some-other-name-1\"\nsome-content-1\n@@ name=\"some-other-name-2\"\nsome-content-2",
				},
			},
			wantErr: nil,
		},
		{
			name: "Should not accept file if the name prefix is not the expected one, even though the prefix is in the name",
			mocks: mocks{
				nextCollectResult: []error{
					nil,
				},
			},
			args: data{
				path:    "some-path/not_mul_something",
				content: "@@ name=\"some-other-name-1\"\nsome-content-1\n@@ name=\"some-other-name-2\"\nsome-content-2",
			},
			want: []data{
				{
					path:    "some-path/not_mul_something",
					content: "@@ name=\"some-other-name-1\"\nsome-content-1\n@@ name=\"some-other-name-2\"\nsome-content-2",
				},
			},
			wantErr: nil,
		},
		{
			name: "Should propagate error if next collector throws one at first file",
			mocks: mocks{
				nextCollectResult: []error{
					errors.New("some-first-file-error"),
				},
			},
			args: data{
				path:    "some-path/mul_something",
				content: "@@ name=\"some-other-name-1\"\nsome-content-1\n@@ name=\"some-other-name-2\"\nsome-content-2",
			},
			wantErr: errors.New("some-first-file-error"),
		},
		{
			name: "Should propagate error if next collector throws one at second file",
			mocks: mocks{
				nextCollectResult: []error{
					nil,
					errors.New("some-second-file-error"),
				},
			},
			args: data{
				path:    "some-path/mul_something",
				content: "@@ name=\"some-other-name-1\"\nsome-content-1\n@@ name=\"some-other-name-2\"\nsome-content-2",
			},
			want: []data{
				{
					path:    "some-other-name-1",
					content: "some-content-1\n",
				},
			},
			wantErr: errors.New("some-second-file-error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mc := &mockCollector{}
			for _, err := range tt.mocks.nextCollectResult {
				mc.On("Collect", mock.Anything).Return(err).Once()
			}
			mc.On("OnPipelineCompleted").Return(nil)
			p := NewSplitterCollector(mc)

			err := p.Collect(&pipeline.Template{
				Path:   tt.args.path,
				Reader: io.NopCloser(strings.NewReader(tt.args.content)),
			})

			assertutils.AssertEqualErrors(t, tt.wantErr, err)
			calls := mc.Calls
			assert.Equal(t, len(tt.mocks.nextCollectResult), len(calls))
			for i := 0; i < len(tt.want); i++ {
				got := calls[i].Arguments.Get(0).(*pipeline.Template)
				assert.Equal(t, tt.want[i].path, got.Path)
				gotContent, err := io.ReadAll(got.Reader)
				assert.NoError(t, err)
				assert.Equal(t, tt.want[i].content, string(gotContent))
			}
		})
	}
}

func Test_splitterCollector_OnPipelineCompleted(t *testing.T) {
	tests := []struct {
		name          string
		nextError     error
		expectedError error
		nextExists    bool
	}{
		{
			name:          "Should return nil when no next collector exists",
			nextExists:    false,
			expectedError: nil,
		},
		{
			name:          "Should return nil when next collector returns nil",
			nextExists:    true,
			nextError:     nil,
			expectedError: nil,
		},
		{
			name:          "Should return error when next collector returns error",
			nextExists:    true,
			nextError:     errors.New("next-collector-error"),
			expectedError: errors.New("next-collector-error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nextCollector pipeline.Collector
			if tt.nextExists {
				nextCollector = &mockCollector{}
				nextCollector.(*mockCollector).On("OnPipelineCompleted").Return(tt.nextError)
			}

			p := NewSplitterCollector(nextCollector)

			err := p.OnPipelineCompleted()

			assertutils.AssertEqualErrors(t, tt.expectedError, err)
			if tt.nextExists {
				nextCollector.(*mockCollector).AssertCalled(t, "OnPipelineCompleted")
			}
		})
	}
}

func Test_splitterCollector_CreateHeaderWithName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should return correct header",
			args: args{
				name: "some-name",
			},
			want: "@@ name=\"some-name\"",
		},
		{
			name: "Should return correct header with a different name",
			args: args{
				name: "some-other-name",
			},
			want: "@@ name=\"some-other-name\"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewSplitterCollector(nil)

			got := p.CreateHeaderWithName(tt.args.name)

			assert.Equal(t, tt.want, got)
		})
	}
}
