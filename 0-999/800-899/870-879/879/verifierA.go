package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseA struct {
	n   int
	pd  [][2]int
	exp string
}

func solveA(n int, pd [][2]int) string {
	last := 0
	for i := 0; i < n; i++ {
		s := pd[i][0]
		d := pd[i][1]
		day := s
		for day <= last {
			day += d
		}
		last = day
	}
	return fmt.Sprint(last)
}

func generateTests() []testCaseA {
	rng := rand.New(rand.NewSource(1))
	cases := make([]testCaseA, 100)
	for i := range cases {
		n := rng.Intn(10) + 1
		pd := make([][2]int, n)
		for j := range pd {
			s := rng.Intn(100) + 1
			d := rng.Intn(50) + 1
			pd[j] = [2]int{s, d}
		}
		cases[i] = testCaseA{n: n, pd: pd, exp: solveA(n, pd)}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", tc.n)
		for _, p := range tc.pd {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
