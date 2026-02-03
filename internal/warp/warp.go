package warp

import (
	"context"
	"errors"
	"os"

	"github.com/waduhek/tmux-warp/internal/parser"
)

// ErrNameNotFound is returned when the requested name does not exist in the
// config file.
var ErrNameNotFound = errors.New("provided key does not exist")

type Warper interface {
	// GetWarpPathByName gets the path mapped against the provided name in the
	// provided config file path. Returns the error [ErrKeyNotFound] if the
	// provided key was not found in the file or an error of the type
	// [*os.PathError] if an error occurred while reading the file.
	GetWarpPathByName(
		ctx context.Context,
		configPath string,
		name string,
	) (string, error)
}

// NewWarper creates a new instance of a warper.
func NewWarper(p parser.Parser) *warper {
	return &warper{p}
}

type warper struct {
	p parser.Parser
}

func (w *warper) GetWarpPathByName(
	ctx context.Context,
	configPath string,
	name string,
) (string, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	resultCh := make(chan parser.ParserResult)
	go w.p.ParseWarpRC(ctx, configPath, resultCh)

	for result := range resultCh {
		// TODO: Print the error to stderr if the result contains an error.
		if result.Error == nil {
			output := result.Output

			if output.Name == name {
				return output.Path, nil
			}
		} else if _, ok := result.Error.(*os.PathError); ok {
			// os.PathError is only returned when an error has occurred while
			// reading the file. Continuing the loop in this case is pointless.
			return "", result.Error
		}
	}

	return "", ErrNameNotFound
}
