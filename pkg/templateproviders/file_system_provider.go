package templateproviders

import (
	"os"
	"path/filepath"

	"github.com/go-scaffold/go-sdk/v2/pkg/filters"
	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
	"github.com/pasdam/files-index/pkg/filesindex"
)

var open = os.Open

type fileSystemProvider struct {
	filter  filters.Filter
	indexer *filesindex.Indexer
}

// NewFileSystemProvider creates a new instance of a FileProvider that reads
// file from the filesystem.
func NewFileSystemProvider(inputDir string, filter filters.Filter) pipeline.TemplateProvider {
	return &fileSystemProvider{
		filter: filter,
		indexer: &filesindex.Indexer{
			Dir: inputDir,
		},
	}
}

func (p *fileSystemProvider) NextTemplate() (*pipeline.Template, error) {
	for {
		item, err := p.indexer.NextFile()
		if err != nil {
			return nil, err
		}

		absolutePath := filepath.Join(p.indexer.Dir, item.Path())
		relativePath := item.Path()
		if p.filter == nil || p.filter.Accept(relativePath) {
			reader, err := open(absolutePath)
			if err != nil {
				return nil, err
			}

			return &pipeline.Template{
				Reader: reader,
				Path:   relativePath,
			}, nil
		}
	}
}
