package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type testCaseB struct {
	n int
	a int
	b int
}

func parseCases(path string) ([]testCaseB, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseB, T)
	for i := 0; i < T; i++ {
		var n, a, b int
		if _, err := fmt.Fscan(in, &n, &a, &b); err != nil {
			// maybe a and b each on separate lines, scanning automatically handles
		}
		// Wait: our file has n\n a \n b each on separate lines; Fscan will read sequential tokens and is ok.
		cases[i] = testCaseB{n: n, a: a, b: b}
	}
	return cases, nil
}

func canPack(pieces []int, n, m int) bool {
	bins := make([]int, m)
	for i := range bins {
		bins[i] = n
	}
	var dfs func(int) bool
	dfs = func(idx int) bool {
		if idx == len(pieces) {
			return true
		}
		p := pieces[idx]
		for i := 0; i < m; i++ {
			if bins[i] >= p {
				bins[i] -= p
				if dfs(idx + 1) {
					return true
				}
				bins[i] += p
			}
		}
		return false
	}
	return dfs(0)
}

func solve(tc testCaseB) int {
	pieces := []int{tc.a, tc.a, tc.a, tc.a, tc.b, tc.b}
	sort.Slice(pieces, func(i, j int) bool { return pieces[i] > pieces[j] })
	total := 4*tc.a + 2*tc.b
	minBars := (total + tc.n - 1) / tc.n
	for m := minBars; m <= 6; m++ {
		if canPack(pieces, tc.n, m) {
			return m
		}
	}
	return 6
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesB.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d\n%d\n%d\n", tc.n, tc.a, tc.b)
		expected := solve(tc)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(gotStr))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
