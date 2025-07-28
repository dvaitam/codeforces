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
	expected string
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(99) + 2 // 2..100
	var in strings.Builder
	var out strings.Builder
	in.WriteString("1\n")
	in.WriteString(fmt.Sprintf("%d\n", n))

	perm := make([]int, n)
	for i := 0; i < n; i++ {
		perm[i] = i + 1
	}
	for i := 1; i < n; i += 2 {
		perm[i], perm[i-1] = perm[i-1], perm[i]
	}
	if n%2 == 1 {
		perm[n-1], perm[n-2] = perm[n-2], perm[n-1]
	}
	for i, v := range perm {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", v)
	}
	out.WriteByte('\n')

	return testCase{input: in.String(), expected: out.String()}
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
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// basic fixed case
	cases := []testCase{{input: "1\n2\n", expected: "2 1\n"}}
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
