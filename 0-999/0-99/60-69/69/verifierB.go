package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n, m       int
	l, r, t, c []int
	expected   int
}

func solveCase(n, m int, l, r, t, c []int) int {
	total := 0
	for j := 1; j <= n; j++ {
		bestTime := 1<<31 - 1
		bestIdx := -1
		for i := 0; i < m; i++ {
			if l[i] <= j && j <= r[i] {
				if t[i] < bestTime || (t[i] == bestTime && i < bestIdx) {
					bestTime = t[i]
					bestIdx = i
				}
			}
		}
		if bestIdx != -1 {
			total += c[bestIdx]
		}
	}
	return total
}

func generateTests() []testCaseB {
	rng := rand.New(rand.NewSource(2))
	cases := make([]testCaseB, 100)
	for idx := range cases {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		l := make([]int, m)
		r := make([]int, m)
		tVals := make([]int, m)
		cVals := make([]int, m)
		for i := 0; i < m; i++ {
			l[i] = rng.Intn(n) + 1
			r[i] = l[i] + rng.Intn(n-l[i]+1)
			tVals[i] = rng.Intn(1000) + 1
			cVals[i] = rng.Intn(1000) + 1
		}
		exp := solveCase(n, m, l, r, tVals, cVals)
		cases[idx] = testCaseB{n: n, m: m, l: l, r: r, t: tVals, c: cVals, expected: exp}
	}
	return cases
}

func run(bin string, input string) (string, error) {
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for j := 0; j < tc.m; j++ {
			fmt.Fprintf(&sb, "%d %d %d %d\n", tc.l[j], tc.r[j], tc.t[j], tc.c[j])
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != fmt.Sprint(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, tc.expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
