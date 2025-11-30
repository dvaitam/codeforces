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
	"0 1 2 3 4 5 6 7 8 9",
	"5 1 1 1 1 1 1 1 1 1",
	"1 10 9 8 7 6 5 4 3 2",
	"10 3 4 5 6 7 8 9 10 11",
	"9 1 2 3 4 5 6 7 8 9",
	"7 4 4 4 4 4 4 4 4 4",
	"15 3 8 9 6 7 5 4 3 2",
	"0 9 9 9 9 9 9 9 9 9",
	"8 1 1 2 2 2 2 2 2 2",
	"20 2 2 2 2 2 2 2 2 2",
}

func referenceSolve(v int, cost []int) string {
	minCost := cost[0]
	for _, c := range cost[1:] {
		if c < minCost {
			minCost = c
		}
	}
	if v < minCost {
		return "-1"
	}
	length := v / minCost
	rem := v - length*minCost
	res := make([]byte, length)
	for i := 0; i < length; i++ {
		for d := 8; d >= 0; d-- {
			extra := cost[d] - minCost
			if extra <= rem {
				res[i] = byte('1' + d)
				rem -= extra
				break
			}
		}
	}
	return string(res)
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

func parseLine(line string) (int, []int, error) {
	parts := strings.Fields(line)
	if len(parts) != 10 {
		return 0, nil, fmt.Errorf("expected 10 numbers, got %d", len(parts))
	}
	v, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, nil, fmt.Errorf("parse v: %w", err)
	}
	cost := make([]int, 9)
	for i := 0; i < 9; i++ {
		c, err := strconv.Atoi(parts[i+1])
		if err != nil {
			return 0, nil, fmt.Errorf("parse cost %d: %w", i+1, err)
		}
		cost[i] = c
	}
	return v, cost, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/binary")
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
		v, cost, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := fmt.Sprintf("%d\n%s\n", v, strings.Join(strings.Fields(line)[1:], " "))
		expected := referenceSolve(v, cost)
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
