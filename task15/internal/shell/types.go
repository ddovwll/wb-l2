package shell

import (
	"context"
	"io"
	"sync"
)

type builtinFunc func(ctx context.Context, args []string, in io.Reader, out io.Writer) error

type proc struct {
	name       string
	args       []string
	isBuiltin  bool
	inputFile  string
	outputFile string
	appendOut  bool
}

type Shell struct {
	state *shellState
}

type shellState struct {
	mu     sync.Mutex
	pgid   int
	cancel context.CancelFunc
}
