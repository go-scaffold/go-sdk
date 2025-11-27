package collectors

import (
	"errors"
	"testing"

	"github.com/go-scaffold/go-sdk/v2/pkg/filters"
	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-utils/pkg/assertutils"
	"github.com/stretchr/testify/assert"
)

func TestNewFilterCollector(t *testing.T) {
	type args struct {
		filter        filters.Filter
		nextCollector pipeline.Collector
	}
	tests := []struct {
		name string
		args args
		want pipeline.Collector
	}{
		{
			name: "Should create collector without next one",
			args: args{
				filter: filters.NewNoOpFilter(),
			},
		},
		{
			name: "Should create collector with next one",
			args: args{
				filter:        filters.NewNoOpFilter(),
				nextCollector: &fileWriterCollector{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFilterCollector(tt.args.filter, tt.args.nextCollector).(*filterCollector)

			assert.Equal(t, tt.args.filter, got.filter)
			assert.Equal(t, tt.args.nextCollector, got.next)
		})
	}
}

func Test_filterCollector_Collect(t *testing.T) {
	filter, err := filters.NewPatternFilter(true, "some-matching-pattern")
	assert.NoError(t, err)

	type fields struct {
		next   pipeline.Collector
		filter filters.Filter
	}
	type args struct {
		path string
	}
	tests := []struct {
		name            string
		args            args
		wantErr         error
		wantNextCollect bool
	}{
		{
			name: "Should collect template if path is accepted by the filter and return no error if the next collector doesn't throw one",
			args: args{
				path: "some-matching-pattern",
			},
			wantNextCollect: true,
		},
		{
			name: "Should collect template if path is accepted by the filter and propagate the error if the next collector throws one",
			args: args{
				path: "some-matching-pattern",
			},
			wantErr:         errors.New("some-next-collector-error"),
			wantNextCollect: true,
		},
		{
			name: "Should not collect template if path is not accepted by the filter",
			args: args{
				path: "some-not-matching-pattern",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tpl := &pipeline.Template{
				Path: tt.args.path,
			}
			collector := &mockCollector{}
			collector.On("Collect", tpl).Return(tt.wantErr)
			collector.On("OnPipelineCompleted").Return(nil) // Also expect the new method to be called
			p := &filterCollector{
				baseCollector: baseCollector{
					next: collector,
				},
				filter: filter,
			}

			err := p.Collect(tpl)

			assertutils.AssertEqualErrors(t, tt.wantErr, err)
			if tt.wantNextCollect {
				collector.AssertCalled(t, "Collect", tpl)
			}
		})
	}
}

func Test_filterCollector_OnPipelineCompleted(t *testing.T) {
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

			p := &filterCollector{
				baseCollector: baseCollector{
					next: nextCollector,
				},
			}

			err := p.OnPipelineCompleted()

			assertutils.AssertEqualErrors(t, tt.expectedError, err)
			if tt.nextExists {
				nextCollector.(*mockCollector).AssertCalled(t, "OnPipelineCompleted")
			}
		})
	}
}
