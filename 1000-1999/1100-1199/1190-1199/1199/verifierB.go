package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseB struct {
	h, l int
}

func parseTestsB(path string) ([]testCaseB, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	cases := make([]testCaseB, 0, len(lines))
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		parts := strings.Fields(ln)
		if len(parts) != 2 {
			return nil, fmt.Errorf("bad line: %q", ln)
		}
		h, _ := strconv.Atoi(parts[0])
		l, _ := strconv.Atoi(parts[1])
		cases = append(cases, testCaseB{h: h, l: l})
	}
	return cases, nil
}

func solveB(tc testCaseB) float64 {
	h := float64(tc.h)
	l := float64(tc.l)
	return (l*l - h*h) / (2 * h)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestsB("testcasesB.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.h, tc.l)
		expect := solveB(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(got), 64)
		if err != nil {
			fmt.Printf("case %d: cannot parse output %q\n", i+1, got)
			os.Exit(1)
		}
		diff := math.Abs(val - expect)
		if diff > 1e-6*math.Max(1.0, math.Abs(expect)) {
			fmt.Printf("case %d failed: expected %.10f got %.10f\n", i+1, expect, val)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
