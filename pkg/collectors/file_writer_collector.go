package collectors

import (
	"bytes"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/go-io-utilx/pkg/ioutilx"
)

type FileWriterCollectorOptions struct {
	OutDir           string
	SkipUnchanged    bool
	CleanupUntracked bool // Flag to control whether to remove untracked files; defaults to false to maintain current behavior
}

type fileWriterCollector struct {
	baseCollector

	opts           FileWriterCollectorOptions
	generatedFiles map[string]bool // Track files generated during pipeline execution
}

func NewFileWriterCollector(outDir string, nextCollector pipeline.Collector) pipeline.Collector {
	return NewFileWriterCollectorWithOpts(
		FileWriterCollectorOptions{
			OutDir:           outDir,
			SkipUnchanged:    false, // Disabled by default
			CleanupUntracked: false, // Disabled by default to maintain current behavior
		},
		nextCollector,
	)
}

// NewFileWriterCollector creates a file writer collector with the provided options
func NewFileWriterCollectorWithOpts(opts FileWriterCollectorOptions, nextCollector pipeline.Collector) pipeline.Collector {
	return &fileWriterCollector{
		opts:           opts,
		generatedFiles: make(map[string]bool),
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

	p.generatedFiles[outPath] = true

	if p.next == nil {
		return nil
	}

	return p.next.Collect(&pipeline.Template{
		Path:   args.Path,
		Reader: io.NopCloser(bytes.NewReader(contentBytes)),
	})
}

func (p *fileWriterCollector) OnPipelineCompleted() error {
	// If cleanup is enabled, remove files that were not generated during pipeline execution
	if p.opts.CleanupUntracked {
		err := p.cleanupUntrackedFiles()
		if err != nil {
			return err
		}
	}

	if p.next == nil {
		return nil
	}
	return p.next.OnPipelineCompleted()
}

// cleanupUntrackedFiles removes files from the output directory that were not generated during the pipeline execution
func (p *fileWriterCollector) cleanupUntrackedFiles() error {
	// Walk through the output directory
	err := filepath.Walk(p.opts.OutDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// If there's an error accessing a file, continue processing other files
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// If the file was not generated during pipeline execution, remove it
		if !p.generatedFiles[path] {
			slog.Info("Removing untracked file", slog.String("path", path))
			err := os.Remove(path)
			if err != nil {
				// Log the error but continue processing other files
				return nil
			}
		}

		return nil
	})

	return err
}
