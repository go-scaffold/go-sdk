package collectors

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/go-scaffold/go-sdk/v2/pkg/filters"
	"github.com/go-scaffold/go-sdk/v2/pkg/pipeline"
)

var (
	defaultMultiFileTemplateNamePrefix    = "mul_" // TODO: make this configurable
	defaultMultiFileTemplateHeadersPrefix = "@@ "  // TODO: make this configurable
)

type SplitterCollector struct {
	baseCollector

	headerPrefix string
	filter       filters.Filter
}

func NewSplitterCollector(nextCollector pipeline.Collector) *SplitterCollector {
	filter, _ := filters.NewPatternFilter(true, fmt.Sprintf("^%s.*", defaultMultiFileTemplateNamePrefix))

	return &SplitterCollector{
		baseCollector: baseCollector{
			next: nextCollector,
		},
		headerPrefix: defaultMultiFileTemplateHeadersPrefix,
		filter:       filter,
	}
}

func (p *SplitterCollector) Collect(args *pipeline.Template) error {
	if !p.filter.Accept(filepath.Base(args.Path)) {
		return p.next.Collect(args)
	}

	scanner := bufio.NewScanner(args.Reader)
	scanner.Split(scanLines)
	var buffer *bytes.Buffer
	var currentTemplate *pipeline.Template
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, p.headerPrefix) { // is header, this indicates a new file
			if currentTemplate != nil { // collect previous file
				err := p.next.Collect(currentTemplate)
				if err != nil {
					return err
				}
			}
			buffer = &bytes.Buffer{}
			currentTemplate = &pipeline.Template{
				Reader: io.NopCloser(buffer),
				Path:   strings.ReplaceAll(strings.TrimPrefix(strings.TrimSpace(line), "@@ name="), "\"", ""), // TODO
			}
		} else if buffer == nil {
			slog.Error("Invalid first line", slog.String("templatePath", args.Path), slog.String("line", line))
			return fmt.Errorf("invalid first line")

		} else {
			buffer.WriteString(line)
		}
	}

	if currentTemplate != nil {
		err := p.next.Collect(currentTemplate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *SplitterCollector) CreateHeaderWithName(name string) string {
	return fmt.Sprintf("%sname=\"%s\"", defaultMultiFileTemplateHeadersPrefix, name)
}

// scanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func scanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0 : i+1], nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), data, nil
	}
	// Request more data.
	return 0, nil, nil
}

func (p *SplitterCollector) OnPipelineCompleted() error {
	if p.next == nil {
		return nil
	}
	return p.next.OnPipelineCompleted()
}
