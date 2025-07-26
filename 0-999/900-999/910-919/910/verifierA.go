package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseA struct {
	n int
	d int
	s string
}

func parseCases(path string) ([]testCaseA, error) {
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
	cases := make([]testCaseA, T)
	for i := 0; i < T; i++ {
		var n, d int
		if _, err := fmt.Fscan(in, &n, &d); err != nil {
			return nil, err
		}
		var s string
		if _, err := fmt.Fscan(in, &s); err != nil {
			return nil, err
		}
		cases[i] = testCaseA{n: n, d: d, s: s}
	}
	return cases, nil
}

func solve(tc testCaseA) int {
	const inf = int(1e9)
	dist := make([]int, tc.n)
	for i := range dist {
		dist[i] = inf
	}
	dist[0] = 0
	for i := 0; i < tc.n; i++ {
		if dist[i] == inf {
			continue
		}
		for j := i + 1; j <= i+tc.d && j < tc.n; j++ {
			if tc.s[j] == '1' && dist[i]+1 < dist[j] {
				dist[j] = dist[i] + 1
			}
		}
	}
	if dist[tc.n-1] == inf {
		return -1
	}
	return dist[tc.n-1]
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
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseCases("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n%s\n", tc.n, tc.d, tc.s)
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
