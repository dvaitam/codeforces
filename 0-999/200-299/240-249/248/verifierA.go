package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var testcases = []string{
	"7 1 0 1 1 1 1 1 1 0 0 1 0 0 1",
	"9 0 1 0 0 1 1 0 1 1 1 0 1 1 1 0 0 0 1",
	"1 1 1",
	"4 1 0 0 0 0 0 1 0",
	"2 1 1 0 1",
	"9 1 0 1 0 1 1 0 1 1 0 1 0 0 0 0 1 1 0",
	"2 0 0 0 0",
	"9 1 1 0 0 1 1 1 1 1 0 1 0 1 1 0 0 0 1",
	"2 0 1 0 1",
	"7 0 0 0 0 0 0 0 0 0 0 1 0 1 0",
	"1 0 0",
	"3 0 1 0 0 0 1",
	"10 0 1 0 0 0 1 1 1 0 0 1 0 0 1 0 1 1 1 0 0",
	"1 0 0",
	"6 1 0 1 0 0 1 1 1 1 1 1 0",
	"9 0 1 0 1 0 1 0 0 1 1 1 1 0 1 1 1 0 0",
	"10 0 1 0 0 0 1 1 1 0 1 1 0 0 1 0 1 0 1 1 0",
	"1 1 1",
	"5 1 0 1 0 0 0 0 1 1 1",
	"1 0 0",
	"1 0 0",
	"2 0 1 1 0",
	"2 1 1 0 0",
	"5 1 0 1 0 1 0 0 1 0 0",
	"1 0 1",
	"9 1 1 0 1 1 1 1 0 0 1 1 0 0 0 1 1 1 1",
	"2 1 0 0 1",
	"3 0 1 1 1 0 1",
	"2 1 0 0 1",
	"3 0 1 1 1 1 1",
	"2 0 1 1 1",
	"6 0 1 0 1 1 0 1 1 0 0 1 0",
	"2 0 0 0 0",
	"7 0 0 1 1 1 1 0 1 0 1 0 1 0 1",
	"4 1 0 0 0 1 0 0 1",
	"7 1 0 0 0 0 0 0 1 1 1 0 0 1 0",
	"10 1 1 1 1 1 1 0 0 0 1 1 1 0 0 1 0 1 1 1 0",
	"2 0 0 0 0",
	"1 1 0",
	"8 1 0 0 1 1 1 0 0 0 1 0 1 1 1 1 0",
	"1 0 0",
	"5 1 1 0 1 1 1 1 1 1 0",
	"3 0 0 1 1 0 0",
	"8 1 0 1 0 0 1 1 0 1 1 1 0 0 1 0 0",
	"1 0 1",
	"2 1 0 1 1",
	"2 0 1 1 0",
	"10 1 0 1 1 0 0 0 1 1 1 0 1 1 0 0 0 1 1 0 0",
	"4 0 1 1 0 0 1 1 0",
	"2 1 0 0 1",
	"3 0 1 1 1 0 1",
	"6 1 0 1 0 0 1 0 1 1 0 0 0",
	"6 0 0 0 0 0 0 0 1 1 0 0 0",
	"3 1 0 1 1 1 0",
	"1 0 1",
	"6 1 1 1 1 0 0 0 1 0 1 0 0",
	"4 1 0 0 0 1 1 1 0",
	"3 0 1 0 0 0 1",
	"5 1 0 1 1 1 0 0 1 1 0",
	"7 1 0 1 1 1 0 1 1 1 0 0 0 1 0",
	"9 0 1 1 0 1 1 0 0 1 1 0 1 1 0 1 0 1 0",
	"8 1 0 0 1 0 1 0 1 1 0 1 0 1 1 0 0",
	"5 0 0 1 1 1 1 0 1 1 1",
	"4 0 0 0 1 1 1 0 0",
	"3 0 1 1 1 0 0",
	"9 0 1 1 1 1 1 1 0 0 0 0 1 1 1 0 1 1 1",
	"7 0 1 1 1 1 0 1 1 0 1 0 0 0 1",
	"8 1 0 0 0 1 0 1 0 0 0 1 1 0 0 1 0",
	"2 1 1 1 1",
	"3 1 1 1 0 1 0",
	"9 1 0 0 0 1 1 1 1 1 1 0 0 1 0 0 1 0 0",
	"2 0 1 0 0",
	"4 1 0 1 1 1 1 0 1",
	"10 1 1 1 1 0 1 1 0 0 0 0 1 0 1 0 0 0 1 0 1",
	"9 0 1 0 1 0 0 1 1 1 1 1 0 0 0 1 0 0 1",
	"4 0 1 1 1 1 1 1 1",
	"4 0 1 0 0 1 1 1 0",
	"7 1 1 1 0 0 0 0 0 0 1 0 1 0 1",
	"9 1 0 0 0 1 0 0 1 1 1 1 0 1 0 1 1 1 1",
	"5 0 0 0 0 0 1 0 0 1 0",
	"9 1 0 0 1 1 0 1 0 1 0 0 0 1 1 0 1 1 1",
	"10 0 0 0 1 1 0 0 1 0 0 0 0 0 1 0 0 0 0 1 1",
	"10 0 1 0 0 1 1 1 1 0 0 0 1 1 1 0 1 0 1 0 1",
	"5 1 0 0 0 1 0 1 1 0 1",
	"6 0 1 0 1 0 0 1 0 1 0 0 0",
	"1 0 0",
	"6 0 0 1 1 0 0 1 1 0 1 1 0",
	"6 1 1 0 1 1 1 0 0 1 1 0 0",
	"9 0 0 0 1 0 0 0 0 0 0 0 1 0 1 0 0 1 1",
	"7 0 0 1 0 1 1 0 1 0 1 0 1 1 1",
	"8 1 0 1 1 1 0 1 0 1 1 0 1 1 1 0 0",
	"5 0 1 1 0 1 1 1 0 0 1",
	"4 1 0 0 1 0 0 0 1",
	"6 1 1 1 1 1 1 0 0 0 0 1 0",
	"6 0 1 1 0 0 1 1 0 0 0 0 0",
	"9 1 1 1 1 1 0 1 1 1 1 1 0 1 0 0 0 1 1",
	"8 0 1 0 0 1 0 0 1 1 0 0 0 1 0 1 0",
	"10 1 0 1 0 0 0 1 0 1 1 1 1 1 1 1 0 1 0 1 0",
	"5 0 1 1 1 1 0 1 1 0 1",
	"6 0 1 1 0 0 0 1 1 0 0 0 0",
}

func referenceSolve(n int, pairs []int) string {
	openLeft, openRight := 0, 0
	for i := 0; i < n; i++ {
		openLeft += pairs[2*i]
		openRight += pairs[2*i+1]
	}
	closeLeft := n - openLeft
	closeRight := n - openRight
	movesLeft := openLeft
	if closeLeft < movesLeft {
		movesLeft = closeLeft
	}
	movesRight := openRight
	if closeRight < movesRight {
		movesRight = closeRight
	}
	return strconv.Itoa(movesLeft + movesRight)
}

func runBinary(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), stderr.String(), err
}

func parseTestCase(line string) (int, []int, error) {
	parts := strings.Fields(line)
	if len(parts) < 1 {
		return 0, nil, fmt.Errorf("empty test case")
	}
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, fmt.Errorf("parse n: %w", err)
	}
	expectedNums := 1 + 2*n
	if len(parts) != expectedNums {
		return 0, nil, fmt.Errorf("expected %d numbers, got %d", expectedNums, len(parts))
	}
	pairs := make([]int, 2*n)
	for i := 0; i < 2*n; i++ {
		v, err := strconv.Atoi(parts[1+i])
		if err != nil {
			return 0, nil, fmt.Errorf("parse value %d: %w", i+1, err)
		}
		pairs[i] = v
	}
	return n, pairs, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierA /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	idx := 0
	for _, tc := range testcases {
		line := strings.TrimSpace(tc)
		if line == "" {
			continue
		}
		idx++
		n, pairs, err := parseTestCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		nums := strings.Fields(line)[1:]
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(nums, " "))
		expected := referenceSolve(n, pairs)
		out, stderr, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx, expected, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
