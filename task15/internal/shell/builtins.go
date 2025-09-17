package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var builtins = map[string]builtinFunc{
	"cd":   builtinCd,
	"pwd":  builtinPwd,
	"echo": builtinEcho,
	"kill": builtinKill,
	"ps":   builtinPs,
}

func builtinCd(_ context.Context, args []string, _ io.Reader, _ io.Writer) error {
	path := ""
	if len(args) == 0 {
		path = os.Getenv("HOME")
		if path == "" {
			return fmt.Errorf("HOME not set")
		}
	} else {
		path = args[0]
	}
	if err := os.Chdir(path); err != nil {
		return err
	}
	return nil
}

func builtinPwd(_ context.Context, _ []string, _ io.Reader, out io.Writer) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, dir)
	return err
}

func builtinEcho(_ context.Context, args []string, _ io.Reader, out io.Writer) error {
	_, err := fmt.Fprintln(out, strings.Join(args, " "))
	return err
}

func builtinKill(_ context.Context, args []string, _ io.Reader, _ io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("kill: pid required")
	}
	pid, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("kill: invalid pid: %v", err)
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	return nil
}

func builtinPs(ctx context.Context, _ []string, _ io.Reader, out io.Writer) error {
	cmd := exec.CommandContext(ctx, "ps", "-ef")
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ps failed: %w", err)
	}
	return nil
}
