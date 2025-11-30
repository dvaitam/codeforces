package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solve(n int) string {
	if n == 5 {
		return strings.TrimSpace(strings.Join([]string{
			"1 2 6",
			"1 3 6",
			"2 4 5",
			"4 5 1",
			"3 4",
			"3 5",
		}, "\n"))
	}

	var sb strings.Builder
	half := n >> 1
	for i := 1; i <= half; i++ {
		sb.WriteString(fmt.Sprintf("%d %d 1\n", i, i+half))
	}
	for i := 1; i+half < n; i++ {
		weight := 2*i - 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", i+half, i+half+1, weight))
	}
	for i := 1; i < half; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", i, i+1))
	}
	sb.WriteString("1 3")
	return strings.TrimSpace(sb.String())
}

var testcases = []int{
	5, 5, 5, 7, 6, 10, 10, 7, 7, 9,
	6, 9, 5, 9, 10, 6, 8, 10, 8, 10,
	9, 7, 9, 8, 9, 7, 5, 5, 7, 8,
	7, 8, 8, 9, 6, 9, 6, 6, 6, 5,
	6, 7, 6, 6, 9, 9, 7, 9, 10, 9,
	6, 8, 8, 10, 9, 7, 9, 7, 7, 8,
	6, 8, 10, 10, 8, 10, 9, 6, 8, 7,
	8, 9, 9, 7, 10, 8, 8, 7, 9, 10,
	9, 10, 8, 8, 10, 6, 7, 10, 6, 9,
	7, 8, 7, 7, 10, 9, 9, 9, 9, 10,
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for idx, n := range testcases {
		expected := solve(n)
		input := fmt.Sprintf("%d\n", n)

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected:\n%s\n got:\n%s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(testcases))
}
