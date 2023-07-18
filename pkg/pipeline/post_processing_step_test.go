package pipeline

import (
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_postProcessingStep_Process(t *testing.T) {
	type fields struct {
		nextProcessor PostProcessor
		processor     PostProcessor
	}
	type args struct {
		args *Template
	}
	type mocks struct {
		processorOut     *Template
		processorErr     error
		nextProcessorOut *Template
		nextProcessorErr error
	}
	tests := []struct {
		name    string
		fields  fields
		mocks   mocks
		args    args
		want    *Template
		wantErr string
	}{
		{
			name: "Should return first processor output if no next step is provided",
			fields: fields{
				nextProcessor: nil,
				processor:     &postProcessorMock{},
			},
			mocks: mocks{
				processorOut: &Template{
					Reader: &readCloserMock{},
					Path:   "some-mock-path",
				},
				processorErr: nil,
			},
			args: args{
				args: &Template{
					Reader: &readCloserMock{},
					Path:   "some-args-path",
				},
			},
			want: &Template{
				Reader: &readCloserMock{},
				Path:   "some-mock-path",
			},
			wantErr: "",
		},
		{
			name: "Should propagate error if first processor throws it",
			fields: fields{
				nextProcessor: nil,
				processor:     &postProcessorMock{},
			},
			mocks: mocks{
				processorOut: nil,
				processorErr: errors.New("some-processing-error"),
			},
			args: args{
				args: &Template{
					Reader: &readCloserMock{},
					Path:   "some-args-path",
				},
			},
			want:    nil,
			wantErr: "some-processing-error",
		},
		{
			name: "Should return next processor output",
			fields: fields{
				nextProcessor: &postProcessorMock{},
				processor:     &postProcessorMock{},
			},
			mocks: mocks{
				processorOut: &Template{
					Reader: &readCloserMock{},
					Path:   "some-processor-path",
				},
				processorErr: nil,
				nextProcessorOut: &Template{
					Reader: &readCloserMock{},
					Path:   "some-next-processor-path",
				},
				nextProcessorErr: nil,
			},
			args: args{
				args: &Template{
					Reader: &readCloserMock{},
					Path:   "some-args-path",
				},
			},
			want: &Template{
				Reader: &readCloserMock{},
				Path:   "some-next-processor-path",
			},
			wantErr: "",
		},
		{
			name: "Should propagate error if next processor throws it",
			fields: fields{
				nextProcessor: &postProcessorMock{},
				processor:     &postProcessorMock{},
			},
			mocks: mocks{
				processorOut: &Template{
					Reader: &readCloserMock{},
					Path:   "some-processor-path",
				},
				processorErr:     nil,
				nextProcessorOut: nil,
				nextProcessorErr: errors.New("some-next-processor-error"),
			},
			args: args{
				args: &Template{
					Reader: &readCloserMock{},
					Path:   "some-args-path",
				},
			},
			want: &Template{
				Reader: &readCloserMock{},
				Path:   "some-next-processor-path",
			},
			wantErr: "some-next-processor-error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &postProcessingStep{
				processor: tt.fields.processor,
			}
			tt.fields.processor.(*postProcessorMock).On("Process", tt.args.args).Return(tt.mocks.processorOut, tt.mocks.processorErr)
			if tt.fields.nextProcessor != nil {
				p.nextStep = &postProcessingStep{
					processor: tt.fields.nextProcessor,
				}
				tt.fields.nextProcessor.(*postProcessorMock).On("Process", tt.mocks.processorOut).Return(tt.mocks.nextProcessorOut, tt.mocks.nextProcessorErr)
				if tt.mocks.processorOut != nil {
					tt.mocks.processorOut.Reader.(*readCloserMock).On("Close").Return(nil)
				}
			}
			tt.args.args.Reader.(*readCloserMock).On("Close").Return(nil)

			got, err := p.Process(tt.args.args)

			tt.args.args.Reader.(*readCloserMock).AssertCalled(t, "Close")
			if tt.wantErr == "" {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)

			} else {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}

type readCloserMock struct {
	io.ReadCloser
	mock.Mock
}

func (m *readCloserMock) Read(p []byte) (n int, err error) {
	res := m.Called(p)
	return res.Int(0), res.Error(1)
}

func (m *readCloserMock) Close() error {
	res := m.Called()
	return res.Error(0)
}
