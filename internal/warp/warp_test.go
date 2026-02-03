package warp_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/waduhek/tmux-warp/internal/parser"
	"github.com/waduhek/tmux-warp/internal/warp"
)

var warp1Path = "/path/to/warp1"
var warp2Path = "/path/to/warp2"

type defaultParser struct {
}

func (p *defaultParser) ParseWarpRC(
	ctx context.Context,
	path string,
	ch chan<- parser.ParserResult,
) {
	defer close(ch)

	ch <- parser.ParserResult{
		Error: nil,
		Output: &parser.ParserOutput{
			Name: "warp1",
			Path: warp1Path,
		},
	}

	ch <- parser.ParserResult{
		Error: nil,
		Output: &parser.ParserOutput{
			Name: "warp2",
			Path: warp2Path,
		},
	}
}

type errorParser struct {
}

func (p *errorParser) ParseWarpRC(
	ctx context.Context,
	path string,
	ch chan<- parser.ParserResult,
) {
	defer close(ch)

	ch <- parser.ParserResult{
		Error: &os.PathError{
			Op:   "",
			Path: "",
			Err:  errors.New("test err"),
		},
		Output: nil,
	}
}

func TestGetWarpPathByName(t *testing.T) {
	w := warp.NewWarper(&defaultParser{})

	path, err := w.GetWarpPathByName(context.Background(), "", "warp2")
	if err != nil {
		t.Fatalf("error while getting warp path: %s", err)
	}

	if path != warp2Path {
		t.Fatalf("expected path %s but got %s", warp2Path, path)
	}
}

func TestGetWarpPathByName_NameNotFound(t *testing.T) {
	w := warp.NewWarper(&defaultParser{})

	_, err := w.GetWarpPathByName(context.Background(), "", "warp3")
	if err == nil {
		t.Fatal("expected error while getting unknown warp path")
	}
	if !errors.Is(err, warp.ErrNameNotFound) {
		t.Fatalf("expected ErrNameNotFound but got: %s", err)
	}
}

func TestGetWarpPathByName_FileError(t *testing.T) {
	w := warp.NewWarper(&errorParser{})

	_, err := w.GetWarpPathByName(context.Background(), "", "warp1")
	if err == nil {
		t.Fatal("expected error while getting unknown warp path")
	}
	if _, ok := err.(*os.PathError); !ok {
		t.Fatalf("expected os.PathError but got: %s", err)
	}
}
