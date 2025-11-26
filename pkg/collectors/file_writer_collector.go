package collectors

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-io-utilx/pkg/ioutilx"
)

type FileWriterCollectorOptions struct {
	OutDir        string
	SkipUnchanged bool
}

type fileWriterCollector struct {
	baseCollector

	opts FileWriterCollectorOptions
}

func NewFileWriterCollector(outDir string, nextCollector pipeline.Collector) pipeline.Collector {
	return &fileWriterCollector{
		opts: FileWriterCollectorOptions{
			OutDir:        outDir,
			SkipUnchanged: false, // Disabled by default
		},
		baseCollector: baseCollector{
			next: nextCollector,
		},
	}
}

// NewFileWriterCollector creates a file writer collector with the provided options
func NewFileWriterCollectorWithOpts(opts FileWriterCollectorOptions, nextCollector pipeline.Collector) pipeline.Collector {
	return &fileWriterCollector{
		opts: opts,
		baseCollector: baseCollector{
			next: nextCollector,
		},
	}
}

func (p *fileWriterCollector) Collect(args *pipeline.Template) error {
	outPath := filepath.Join(p.opts.OutDir, args.Path)

	contentBytes, err := io.ReadAll(args.Reader)
	if err != nil {
		return err
	}

	writeFile := true
	if p.opts.SkipUnchanged {
		existingContent, err := os.ReadFile(outPath)
		if err == nil && bytes.Equal(existingContent, contentBytes) {
			writeFile = false
		}
	}

	if writeFile {
		err = ioutilx.ReaderToFile(bytes.NewReader(contentBytes), outPath)
		if err != nil {
			return err
		}
	}

	if p.next == nil {
		return nil
	}

	return p.next.Collect(&pipeline.Template{
		Path:   args.Path,
		Reader: io.NopCloser(bytes.NewReader(contentBytes)),
	})
}
