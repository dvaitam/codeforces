package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	a, b, c, d int64
}

func solve(tc testCase) int {
	ad := tc.a * tc.d
	bc := tc.b * tc.c
	if ad == bc {
		return 0
	}
	if ad == 0 || bc == 0 {
		return 1
	}
	if ad%bc == 0 || bc%ad == 0 {
		return 1
	}
	return 2
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	return fmt.Sprintf("1\n%d %d %d %d\n", tc.a, tc.b, tc.c, tc.d)
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{0, 1, 0, 1},
		{1, 1, 2, 2},
		{1, 2, 3, 4},
		{0, 5, 7, 8},
	}
	for len(cases) < 120 {
		a := rng.Int63n(1_000_000_001)
		b := rng.Int63n(1_000_000_000) + 1
		c := rng.Int63n(1_000_000_001)
		d := rng.Int63n(1_000_000_000) + 1
		cases = append(cases, testCase{a: a, b: b, c: c, d: d})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		input := buildInput(tc)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		exp := strconv.Itoa(solve(tc))
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
