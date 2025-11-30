package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// The original testcases file contained 100 blank lines; embed that directly.
const rawTestcases = "\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n"

func parseTestcases() []string {
	// rawTestcases has 100 newline-separated empty strings.
	lines := strings.Split(strings.TrimRight(rawTestcases, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return make([]string, 0)
	}
	return lines
}

// solveCase is the embedded logic from 241G.go, adapted to return the output as a string.
func solveCase() string {
	var sb strings.Builder
	sb.WriteString("302 ")
	sb.WriteByte('\n')
	sb.WriteString(" 0 800000")
	sb.WriteByte('\n')
	s := 60000
	for j := 300; j > 0; j-- {
		sb.WriteString(fmt.Sprintf("%d %d\n", s, j))
		s += 2*j - 1
	}
	sb.WriteString(fmt.Sprintf("%d %d\n", s+(1<<17), 800000))
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	expected := strings.TrimSpace(solveCase())
	testcases := parseTestcases()
	// If parsing trimmed away empties, default to 100 runs to match original file length.
	runCount := len(testcases)
	if runCount == 0 {
		runCount = 100
	}

	for idx := 0; idx < runCount; idx++ {
		cmd := exec.Command(bin)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", runCount)
}
