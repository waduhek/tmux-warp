package parser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
)

// ErrInvalidLineFormat is returned when a line in the file does not match the
// expected format.
var ErrInvalidLineFormat = errors.New("invalid line format")

// ErrContextDone is returned when the provided context is closed while parsing
// the file.
var ErrContextDone = errors.New("context was closed while parsing file")

// ParserOutput represents a parsed row from the config file.
type ParserOutput struct {
	Name string
	Path string
}

// ParserResult is the result from the [ParseWarpRC] function.
type ParserResult struct {
	Output *ParserOutput
	Error  error
}

type Parser interface {
	// ParseWarpRC parses the warprc file at the provided path. The function
	// assumes that the file is formatted in the correct warprc format, i.e.
	// <name>:<path>.
	ParseWarpRC(ctx context.Context, path string, ch chan<- ParserResult)
}

// NewParser creates a new instance of a parser.
func NewParser() *parser {
	return &parser{}
}

type parser struct {
}

// ParseWarpRC parses the warprc file at the provided path. The function assumes
// that the file is formatted in the correct warprc format, i.e. <name>:<path>.
func (p *parser) ParseWarpRC(
	ctx context.Context,
	path string,
	ch chan<- ParserResult,
) {
	defer close(ch)

	file, err := p.getFileHandler(path)
	if err != nil {
		ch <- *err
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	p.scanAndParse(ctx, scanner, ch)
}

func (p *parser) getFileHandler(path string) (*os.File, *ParserResult) {
	file, err := os.Open(path)
	if err != nil {
		result := &ParserResult{
			Output: nil,
			Error:  fmt.Errorf("error while opening file for parsing: %w", err),
		}

		return nil, result
	}

	return file, nil
}

func (p *parser) scanAndParse(
	ctx context.Context,
	scanner *bufio.Scanner,
	ch chan<- ParserResult,
) {
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			result := ParserResult{
				Output: nil,
				Error:  ErrContextDone,
			}

			ch <- result
			return

		default:
			line := scanner.Text()
			p.parseLine(line, ch)
		}
	}
}

func (p *parser) parseLine(line string, ch chan<- ParserResult) {
	lineParts := strings.Split(line, ":")

	var result ParserResult
	if len(lineParts) != 2 {
		result = ParserResult{
			Output: nil,
			Error:  ErrInvalidLineFormat,
		}
	} else {
		result = ParserResult{
			Output: &ParserOutput{
				Name: lineParts[0],
				Path: lineParts[1],
			},
			Error: nil,
		}
	}

	ch <- result
}
