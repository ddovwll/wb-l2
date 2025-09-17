package shell

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

func NewShell() *Shell {
	return &Shell{state: &shellState{}}
}

func (s *Shell) Close() {
	s.state.mu.Lock()
	cancel := s.state.cancel
	pgid := s.state.pgid
	s.state.pgid = 0
	s.state.cancel = nil
	s.state.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	if pgid != 0 {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)

		done := make(chan struct{})
		go func() {
			time.Sleep(300 * time.Millisecond)
			_ = syscall.Kill(-pgid, syscall.SIGKILL)
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(1 * time.Second):
		}
	}
}

func (s *Shell) ExecuteLine(line string) error {
	cmds, ops := splitConditional(line)
	if len(cmds) == 0 {
		return nil
	}
	prevExit := 0
	for i, cmd := range cmds {
		if i > 0 {
			op := ops[i-1]
			if op == "&&" && prevExit != 0 {
				continue
			}
			if op == "||" && prevExit == 0 {
				continue
			}
		}

		exit, err := s.executePipeline(cmd)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		prevExit = exit
	}
	return nil
}

func IsInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (fi.Mode() & os.ModeCharDevice) != 0
}

func (s *Shell) HandleSigInt() {
	s.state.mu.Lock()
	pgid := s.state.pgid
	cancel := s.state.cancel
	s.state.mu.Unlock()
	println()

	if pgid != 0 {
		_ = syscall.Kill(-pgid, syscall.SIGINT)
		if cancel != nil {
			cancel()
		}
	} else {
		curDir, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		fmt.Printf("%s > ", curDir)
	}
}
