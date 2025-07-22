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

type testCase struct {
	input    string
	expected int
}

func compute(a, b []int) int {
	n := len(a)
	cnt := 0
	for i := 0; i < n; i++ {
		open := false
		for j := 0; j < n; j++ {
			if i != j && b[j] == a[i] {
				open = true
				break
			}
		}
		if !open {
			cnt++
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(100) + 1
	a := make([]int, n)
	b := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(1000) + 1
		b[i] = rng.Intn(1000) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", a[i], b[i]))
	}
	return testCase{input: sb.String(), expected: compute(a, b)}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var ans int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &ans); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if ans != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, ans)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []testCase{
		{input: "1\n1 1\n", expected: 1},
		{input: "2\n1 2\n2 1\n", expected: 0},
		{input: "2\n1 2\n2 3\n", expected: 1},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
