package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded testcases (one expression per line).
const embeddedTestcases = `2+1+2+3+2+2+2+2+2+3+1+3+1
1+1+3+2+3+3+3+1+2+1
3+2+2
1+2+2+2+3+3+1+3+2+2+3+2+1+3+1+1+3+2
3
2+1+3+2+3+1+1+3+1+1+1+3+2+1+1+2
2+1+2+3+2+3+1+3+2+3+1+3+3+3+2+2+1
2+2+3+1+2+1+1+1+1+3+3+2+2+1+1+3+1+1+1+1
3+2+3+3+2+3+1+1+3+3+2+3+2+2+2+3+3+3
1+2+3+1+2+3+3+2+1+1+1+3
1+3+1+2+1+2+2+1+1
3+1+1+3+3
3+3+1+1+1+3+1+3+3+1+2+1+2+1+1+3+1+1
3+1+2+1+3+1
3
3+1+2+1+1+1+3+2+2+2+1+1+3+2
3+1
1+2+2+3+2+3+1+3+3+1+1+3+1
2+3+2+1+3+2
1+2+3+2+3+3
3+2+2+3+2+1+3+3+1+2
2+3+1
2+1+1+2+2+3+2+3+2+3+3+3+1+3+2+2+3+2
1+3+1
1+1+1+3+2+2+3+3+3+2+1
3+3+2+3+3+1+1+2+1+2+3+1+2
2+3+3+1+1+2+2+2+2+1+2+1+3+3+1+3+1
2
2+1+1+1+3+1+3+3+3+1+1+1+3+3
2+2+3+1+1+2+2
1+2+2
2+1+3+3
1+1+2+1+1+1+1+3+2+3+2+2
1+3+3+3+3+2+3+3+2+3+2+2+3+1+1+2+3+2+1
1+2+2+2+2
2+3+1
2+1
3+2+2+2+3
2+1+2+3+1
2+1
3+1+2+2+2+2+2+1+1+3+2+2+2+2+1+2+1
2+1+2+2+3+3+1+1+3+3+2+3+1+1+1+1
1+2+1+1+2+3+3+2
2+3+3+3+1+2+1+2+1+2+3+1+2+1+2
1+3+1+3
3+1+1+2+2+2+1+3+1+1+3+1+1+1+2
2+3+1+1+3+2+1+3+2+3+3+2+3
3+2+2+2+3+3+1+3+3+1+1+2+3+3+2+2
3+1
3+1+2+3+2+3+3+2+1
3+1+1
1+1+2+1+2+2+1+1
2+3+2+3+3+1+3+1+3+3+3+1+3+2+1
3+3+2+2+2+3+3+1+1+3
3
2+3+3+2+2+1
2+2+2+2+2+1+1+3+1+1+2+3+2+1+1+2
1+2+3+3+1+3+3+1+2+2+1+3+2+2
1+1+2+1+1+1+3+3+1+3+2+1+3+3+3
1+2+2+1+1+3+3+2+3+3+2+3
3+3+3+2
2+2+3+3+1
1+1+2+2+2+1+2+2+3+1+1+1+2+2+1+3+3
3+1+1+1+3+3+3+3+2+3+2+1+1+2+3+2+1+1+2
1+2+1
2
3+2+1+2+2+3+2+1+2+1+1+2+3+1
2+1+1+1+2+1+3+1+1+1+1+3
3+1+2+2
3
1+1+2+1+2+2+3+2
1+1+2+2+2
2+3+3+2+1+3+1+1+3+3
1+2+1+1+1+2+3+1+1+1
2+2+3+1+1+3+1+2+1+3
1+2+2
2+3+1+3+2+2+3+2+1+1+2+3+2+3+3+1+2+2
2
2+1+2+2+2+1+1+1+2+1+3
3+1+3+2+2+1+2+2+1+1+2+2+3+1+2+2+3+3+1+2
2+1+3+2+3+2+1
2+1+3+2+3+3+1
3+2+1+3+2+1+3+2+2+3+1+1+2+3
1+2
3+3+3+2+2+1+3+2+2+2+3+1+3
1+1+2
3+2+1+3+1+1+1+2+2+2
1+3+3+3+1+3+1+2+3+3+2+3+2+2+2+2+3
1+1+3+1+2
3+2+1+3+3+2+2+3+2+3+2+1+2
3+3+2+2+3+1+2+2+1+3+2+1+3+1+1+2
2+1+3+1+3+1+3+3+3+2+1+2+1+1+3+1
3+3+2+3+1+1+3+2+1
2+2+3+2
1+2+2+3+3+2+1+2+3+1+3+2+1
1+2+2+2+2+2+2
1+2+1+3+3+1+2+1
3+1+1+2+1+3+3+3
1+3+2+3+1+3
2+2+2+2+1+3+3+3+3+3+3+3+2+3+2+2+2`

func normalize(s string) string {
	return strings.TrimSpace(s)
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return normalize(out.String()), nil
}

// Embedded solver logic from 339A.go.
func solveCase(expr string) string {
	parts := strings.Split(expr, "+")
	sort.Strings(parts)
	return strings.Join(parts, "+")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	lines := strings.Split(strings.TrimSpace(embeddedTestcases), "\n")
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		input := line + "\n"
		expected := solveCase(line)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
