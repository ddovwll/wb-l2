package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"qwe/internal/shell"
)

func main() {
	interactive := shell.IsInteractive()
	s := shell.NewShell()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	defer func() {
		signal.Stop(sigs)
		s.Close()
	}()

	go func() {
		for range sigs {
			s.HandleSigInt()
		}
	}()

	reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(reader)

	for {
		curDir, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
		if interactive {
			fmt.Printf("%s > ", curDir)
		}

		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "scan error:", err)
			}
			signal.Stop(sigs)
			s.Close()
			fmt.Println()
			return
		}
		line := scanner.Text()
		line = shell.TrimSpace(line)
		if line == "" {
			continue
		}
		if shell.IsComment(line) {
			continue
		}
		if err := s.ExecuteLine(line); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		}
	}
}
