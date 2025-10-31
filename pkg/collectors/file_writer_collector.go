package collectors

import (
	"bytes"
	"io"
	"path/filepath"

	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-io-utilx/pkg/ioutilx"
)

type fileWriterCollector struct {
	baseCollector

	outDir string
}

func NewFileWriterCollector(outDir string, nextCollector pipeline.Collector) pipeline.Collector {
	return &fileWriterCollector{
		outDir: outDir,
		baseCollector: baseCollector{
			next: nextCollector,
		},
	}
}

func (p *fileWriterCollector) Collect(args *pipeline.Template) error {
	outPath := filepath.Join(p.outDir, args.Path)

	buf := &bytes.Buffer{}
	teeReader := io.TeeReader(args.Reader, buf)

	err := ioutilx.ReaderToFile(teeReader, outPath)
	if err != nil {
		return err
	}

	if p.next == nil {
		return nil
	}

	return p.next.Collect(&pipeline.Template{
		Path:   args.Path,
		Reader: io.NopCloser(buf),
	})
}
