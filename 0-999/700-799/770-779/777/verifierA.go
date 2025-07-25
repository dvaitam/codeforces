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
	n int64
	x int
}

func simulate(start int, steps int64) int {
	pos := start
	for i := int64(0); i < steps; i++ {
		if i%2 == 0 {
			if pos == 0 {
				pos = 1
			} else if pos == 1 {
				pos = 0
			}
		} else {
			if pos == 1 {
				pos = 2
			} else if pos == 2 {
				pos = 1
			}
		}
	}
	return pos
}

func expected(n int64, x int) int {
	steps := n % 6
	for start := 0; start < 3; start++ {
		if simulate(start, steps) == x {
			return start
		}
	}
	return -1
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
	return fmt.Sprintf("%d\n%d\n", tc.n, tc.x)
}

func genCases() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 120)
	for n := int64(0); n <= 6; n++ {
		for x := 0; x < 3; x++ {
			cases = append(cases, testCase{n: n, x: x})
		}
	}
	for len(cases) < 120 {
		n := rng.Int63n(2_000_000_000) + 1
		x := rng.Intn(3)
		cases = append(cases, testCase{n: n, x: x})
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
		exp := strconv.Itoa(expected(tc.n, tc.x))
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
