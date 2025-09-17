package shell

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
)

func (s *Shell) executePipeline(line string) (int, error) {
	parts := splitPipeline(line)
	if len(parts) == 0 {
		return 0, nil
	}
	procs, err := parsePipelineParts(parts)
	if err != nil {
		return 1, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.state.mu.Lock()
	s.state.cancel = cancel
	s.state.mu.Unlock()
	defer func() {
		s.state.mu.Lock()
		s.state.pgid = 0
		s.state.cancel = nil
		s.state.mu.Unlock()
		cancel()
	}()

	type rw struct {
		r *io.PipeReader
		w *io.PipeWriter
	}
	pipes := make([]rw, 0, len(procs)-1)
	for i := 0; i < len(procs)-1; i++ {
		r, w := io.Pipe()
		pipes = append(pipes, rw{r: r, w: w})
	}

	usedReader := make([]bool, len(pipes))
	usedWriter := make([]bool, len(pipes))

	openedFiles := make([]*os.File, 0)
	defer func() {
		for _, f := range openedFiles {
			_ = f.Close()
		}
	}()

	var wg sync.WaitGroup
	errCh := make(chan error, len(procs))
	lastExitCh := make(chan int, 1)
	lastIdx := len(procs) - 1

	for i, p := range procs {
		var in io.Reader = os.Stdin
		var out io.Writer = os.Stdout

		if p.inputFile != "" {
			f, err := os.Open(p.inputFile)
			if err != nil {
				return 1, fmt.Errorf("open input %s: %w", p.inputFile, err)
			}
			openedFiles = append(openedFiles, f)
			in = f
		} else if i > 0 {
			in = pipes[i-1].r
			usedReader[i-1] = true
		}

		if p.outputFile != "" {
			flags := os.O_CREATE | os.O_WRONLY
			if p.appendOut {
				flags |= os.O_APPEND
			} else {
				flags |= os.O_TRUNC
			}
			f, err := os.OpenFile(p.outputFile, flags, 0644)
			if err != nil {
				return 1, fmt.Errorf("open output %s: %w", p.outputFile, err)
			}
			openedFiles = append(openedFiles, f)
			out = f
		} else if i < len(procs)-1 {
			out = pipes[i].w
			usedWriter[i] = true
		}

		if p.isBuiltin {
			fn := builtins[p.name]
			wg.Add(1)
			go func(fn builtinFunc, args []string, in io.Reader, out io.Writer, name string, idx int) {
				defer wg.Done()
				defer closeIfPipeWriter(out)
				err := fn(ctx, args, in, out)
				if idx == lastIdx {
					if err == nil {
						lastExitCh <- 0
					} else {
						lastExitCh <- 1
					}
				}
				if err != nil {
					select {
					case errCh <- fmt.Errorf("builtin %s: %w", name, err):
					default:
					}
				}
			}(fn, p.args, in, out, p.name, i)
		} else {
			wg.Add(1)
			go func(name string, args []string, in io.Reader, out io.Writer, idx int) {
				defer wg.Done()
				defer closeIfPipeWriter(out)

				cmd := exec.CommandContext(ctx, name, args...)
				cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
				cmd.Stdin = in
				cmd.Stdout = out
				cmd.Stderr = os.Stderr

				if err := cmd.Start(); err != nil {
					if idx == lastIdx {
						lastExitCh <- 127
					}
					select {
					case errCh <- fmt.Errorf("start %s: %w", name, err):
					default:
					}
					return
				}

				s.state.mu.Lock()
				if s.state.pgid == 0 {
					s.state.pgid = cmd.Process.Pid
				}
				s.state.mu.Unlock()

				if err := cmd.Wait(); err != nil {
					exit := 1
					var ee *exec.ExitError
					if errors.As(err, &ee) {
						if ws, ok := ee.Sys().(syscall.WaitStatus); ok {
							exit = ws.ExitStatus()
						}
					}
					if idx == lastIdx {
						lastExitCh <- exit
					}
					select {
					case <-ctx.Done():
					default:
						select {
						case errCh <- fmt.Errorf("wait %s: %w", name, err):
						default:
						}
					}
				} else {
					if idx == lastIdx {
						lastExitCh <- 0
					}
				}
			}(p.name, p.args, in, out, i)
		}
	}

	for i := range pipes {
		if !usedReader[i] {
			_ = pipes[i].r.Close()
		}
		if !usedWriter[i] {
			_ = pipes[i].w.Close()
		}
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		select {
		case exit := <-lastExitCh:
			return exit, nil
		default:
			return 0, nil
		}
	case err := <-errCh:
		cancel()
		select {
		case <-done:
		case <-time.After(200 * time.Millisecond):
		}
		select {
		case exit := <-lastExitCh:
			return exit, err
		default:
			return 1, err
		}
	case exit := <-lastExitCh:
		select {
		case <-done:
			return exit, nil
		case <-time.After(500 * time.Millisecond):
			return exit, nil
		}
	}
}

func closeIfPipeWriter(w io.Writer) {
	if pw, ok := w.(*io.PipeWriter); ok {
		_ = pw.Close()
	}
}
