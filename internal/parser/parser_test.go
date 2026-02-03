package parser_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/waduhek/tmux-warp/internal/parser"
	"github.com/waduhek/tmux-warp/internal/testutils"
)

var ctx = context.Background()
var p = parser.NewParser()

func TestParseWarpRC(t *testing.T) {
	testFile := filepath.Join("testdata", "parse_normal.txt")
	destFile := filepath.Join(t.TempDir(), "parse_normal.txt")

	if err := testutils.CopyFile(testFile, destFile); err != nil {
		t.Fatalf("error while copying test file: %s", err)
	}

	resultChan := make(chan parser.ParserResult)
	go p.ParseWarpRC(ctx, destFile, resultChan)

	for result := range resultChan {
		if result.Error != nil {
			t.Errorf("error while parsing file: %s", result.Error)
		}
	}
}

func TestParseWarpRC_InvalidFileFormat(t *testing.T) {
	testFile := filepath.Join("testdata", "parse_invalid_format.txt")
	destFile := filepath.Join(t.TempDir(), "parse_invalid_format.txt")

	if err := testutils.CopyFile(testFile, destFile); err != nil {
		t.Fatalf("error while copying test file: %s", err)
	}

	resultChan := make(chan parser.ParserResult)
	go p.ParseWarpRC(ctx, destFile, resultChan)

	result := <-resultChan
	if result.Error == nil {
		t.Fatalf("expected error while parsing file")
	}

	if !errors.Is(result.Error, parser.ErrInvalidLineFormat) {
		t.Fatalf("expected ErrInvalidLineFormat but got: %s", result.Error)
	}
}

func TestParseWarpRC_BlankLine(t *testing.T) {
	testFile := filepath.Join("testdata", "parse_blank_line.txt")
	destFile := filepath.Join(t.TempDir(), "parse_blank_line.txt")

	if err := testutils.CopyFile(testFile, destFile); err != nil {
		t.Fatalf("error while copying test file: %s", err)
	}

	resultChan := make(chan parser.ParserResult)
	go p.ParseWarpRC(ctx, destFile, resultChan)

	_ = <-resultChan
	result := <-resultChan
	_ = <-resultChan

	if result.Error == nil {
		t.Fatalf("expected error while parsing file")
	}

	if !errors.Is(result.Error, parser.ErrInvalidLineFormat) {
		t.Fatalf("expected ErrInvalidLineFormat but got: %s", result.Error)
	}
}

func TestParseWarpRC_NonExistentFile(t *testing.T) {
	destFile := filepath.Join(t.TempDir(), "non_existent.txt")

	resultChan := make(chan parser.ParserResult)
	go p.ParseWarpRC(ctx, destFile, resultChan)

	result := <-resultChan
	if result.Error == nil {
		t.Fatalf("expected error while parsing file")
	}
}

func TestParseWarpRC_ContextCancellation(t *testing.T) {
	testFile := filepath.Join("testdata", "parse_large.txt")
	destFile := filepath.Join(t.TempDir(), "parse_large.txt")

	if err := testutils.CopyFile(testFile, destFile); err != nil {
		t.Fatalf("error while copying test file: %s", err)
	}

	resultChan := make(chan parser.ParserResult)

	ctx, cancel := context.WithCancel(ctx)
	go p.ParseWarpRC(ctx, destFile, resultChan)
	cancel()

	for result := range resultChan {
		if result.Error != nil &&
			!errors.Is(result.Error, parser.ErrContextDone) {
			t.Error("expected context done error")
		}
	}
}
