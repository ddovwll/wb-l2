package shell

import (
	"fmt"
	"os"
	"strings"
)

func TrimSpace(s string) string { return strings.TrimSpace(s) }
func IsComment(s string) bool   { return strings.HasPrefix(strings.TrimSpace(s), "#") }

func expandToken(s string) string {
	return os.Expand(s, func(varname string) string {
		return os.Getenv(varname)
	})
}

func splitConditional(line string) (cmds []string, ops []string) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}
	cmds = make([]string, 0)
	ops = make([]string, 0)
	start := 0
	for i := 0; i < len(line); i++ {
		if i+1 < len(line) {
			two := line[i : i+2]
			if two == "&&" || two == "||" {
				cmd := strings.TrimSpace(line[start:i])
				if cmd != "" {
					cmds = append(cmds, cmd)
				}
				ops = append(ops, two)
				i += 1
				start = i + 1
			}
		}
	}
	last := strings.TrimSpace(line[start:])
	if last != "" {
		cmds = append(cmds, last)
	}
	return cmds, ops
}

func splitPipeline(line string) []string {
	raw := strings.Split(line, "|")
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func parsePipelineParts(parts []string) ([]proc, error) {
	procs := make([]proc, 0, len(parts))
	for _, part := range parts {
		tokens := strings.Fields(part)
		if len(tokens) == 0 {
			continue
		}
		p := proc{}
		p.name = expandToken(tokens[0])
		if _, ok := builtins[p.name]; ok {
			p.isBuiltin = true
		}
		for i := 1; i < len(tokens); i++ {
			t := tokens[i]
			if t == "<" {
				if i+1 >= len(tokens) {
					return nil, fmt.Errorf("syntax error: expected filename after '<'")
				}
				p.inputFile = expandToken(tokens[i+1])
				i++
				continue
			}
			if t == ">" || t == ">>" {
				if i+1 >= len(tokens) {
					return nil, fmt.Errorf("syntax error: expected filename after '%s'", t)
				}
				p.outputFile = expandToken(tokens[i+1])
				p.appendOut = t == ">>"
				i++
				continue
			}
			p.args = append(p.args, expandToken(t))
		}
		procs = append(procs, p)
	}
	return procs, nil
}
