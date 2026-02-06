package parser

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrInvalidLineFormat is returned when a line in the file does not match the
// expected format.
var ErrInvalidLineFormat = errors.New("invalid line format")

// ErrContextDone is returned when the provided context is closed while parsing
// the file.
var ErrContextDone = errors.New("context was closed while parsing file")

// ErrHomeEnvNotSet is returned when attempting to expand a path aliased to the
// home directory and the value of HOME was not found/set.
var ErrHomeEnvNotSet = errors.New("value of environment variable HOME not set")

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
		result := p.buildErrorResult(
			fmt.Errorf("error while opening file for parsing: %w", err),
		)

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
			result := p.buildErrorResult(ErrContextDone)

			ch <- *result
			return

		default:
			line := scanner.Text()
			p.parseLine(line, ch)
		}
	}
}

func (p *parser) parseLine(line string, ch chan<- ParserResult) {
	lineParts := strings.Split(line, ":")

	if len(lineParts) != 2 {
		result := p.buildErrorResult(ErrInvalidLineFormat)

		ch <- *result
		return
	}

	name := lineParts[0]
	path := lineParts[1]

	result := p.expandPathIfRequired(name, path)
	ch <- *result
}

func (p *parser) expandPathIfRequired(name, path string) *ParserResult {
	if path[0] == '~' {
		expandedPath, err := p.expandHomeAlias(path)
		if err != nil {
			return p.buildErrorResult(err)
		}

		return p.buildResult(name, expandedPath)
	}

	return p.buildResult(name, path)
}

func (p *parser) expandHomeAlias(path string) (string, error) {
	homeEnv := os.Getenv("HOME")
	if homeEnv == "" {
		return "", ErrHomeEnvNotSet
	}

	return filepath.Join(homeEnv, path[1:]), nil
}

func (p *parser) buildResult(name, path string) *ParserResult {
	return &ParserResult{
		Output: &ParserOutput{
			Name: name,
			Path: path,
		},
		Error: nil,
	}
}

func (p *parser) buildErrorResult(err error) *ParserResult {
	return &ParserResult{
		Output: nil,
		Error:  err,
	}
}
