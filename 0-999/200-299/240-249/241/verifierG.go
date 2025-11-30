package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcaseRuns = 100

func solve() string {
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

	expected := strings.TrimSpace(solve())

	for idx := 0; idx < testcaseRuns; idx++ {
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

	fmt.Printf("All %d tests passed\n", testcaseRuns)
}
