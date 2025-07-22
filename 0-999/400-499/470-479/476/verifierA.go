package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, m int) int {
	kStart := (n + 1) / 2
	t := (kStart + m - 1) / m * m
	if t <= n {
		return t
	}
	return -1
}

type testCase struct {
	n int
	m int
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(10000) + 1
	m := rng.Intn(10) + 1
	return testCase{n, m}
}

func runCase(bin string, tc testCase) error {
	input := fmt.Sprintf("%d %d\n", tc.n, tc.m)
	exp := expected(tc.n, tc.m)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	var got int
	if _, err := fmt.Sscan(out, &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{n: 2, m: 1},
		{n: 2, m: 2},
		{n: 5, m: 2},
		{n: 9, m: 3},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
