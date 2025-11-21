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
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	bin := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		out, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
		if err := validate(tc.input, out); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: %v\ninput:\n%soutput:\n%s", i+1, err, tc.input, out)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(input, output string) error {
	inFields := strings.Fields(input)
	ptr := 0
	readInt := func() (int64, error) {
		if ptr >= len(inFields) {
			return 0, fmt.Errorf("unexpected end of input")
		}
		var x int64
		_, err := fmt.Sscan(inFields[ptr], &x)
		ptr++
		return x, err
	}

	tCases, err := readInt()
	if err != nil {
		return fmt.Errorf("failed to read t: %v", err)
	}
	type instance struct {
		n, a, b int64
	}
	inst := make([]instance, tCases)
	for i := int64(0); i < tCases; i++ {
		n, err := readInt()
		if err != nil {
			return fmt.Errorf("failed to read n for case %d: %v", i+1, err)
		}
		a, err := readInt()
		if err != nil {
			return fmt.Errorf("failed to read a for case %d: %v", i+1, err)
		}
		b, err := readInt()
		if err != nil {
			return fmt.Errorf("failed to read b for case %d: %v", i+1, err)
		}
		inst[i] = instance{n: n, a: a, b: b}
	}

	ansTokens := strings.Fields(output)
	if len(ansTokens) < int(tCases) {
		return fmt.Errorf("expected at least %d answers, got %d", tCases, len(ansTokens))
	}

	for i, tc := range inst {
		exp := solve(tc.n, tc.a, tc.b)
		got := strings.ToLower(ansTokens[i])
		if got != "yes" && got != "no" {
			return fmt.Errorf("case %d: invalid answer token %q", i+1, ansTokens[i])
		}
		if (exp && got != "yes") || (!exp && got != "no") {
			return fmt.Errorf("case %d: expected %v got %q", i+1, exp, ansTokens[i])
		}
	}
	return nil
}

func solve(n, a, b int64) bool {
	// Blue block must be centered for symmetry: requires n-b even.
	if (n-b)%2 != 0 {
		return false
	}
	if a <= b {
		// Hide red inside blue.
		return true
	}
	// Otherwise red must also be centerable so leftover red on both sides matches.
	return (n-a)%2 == 0
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Statement sample.
	tests = append(tests, testCase{input: "7\n5 3 1\n4 1 2\n7 7 4\n8 3 7\n1 1 1\n1000000000 100000000 100000000\n3 2 1\n"})

	// Exhaustive small cases.
	for n := 1; n <= 8; n++ {
		var b strings.Builder
		count := 0
		for a := 1; a <= n; a++ {
			for bb := 1; bb <= n; bb++ {
				count++
				fmt.Fprintf(&b, "%d %d %d\n", n, a, bb)
			}
		}
		tests = append(tests, testCase{input: fmt.Sprintf("%d\n%s", count, b.String())})
	}

	// Random larger values.
	for i := 0; i < 20; i++ {
		t := rng.Intn(20) + 1
		var b strings.Builder
		for j := 0; j < t; j++ {
			n := int64(rng.Intn(1_000_000_000) + 1)
			a := int64(rng.Intn(int(n)) + 1)
			bb := int64(rng.Intn(int(n)) + 1)
			fmt.Fprintf(&b, "%d %d %d\n", n, a, bb)
		}
		tests = append(tests, testCase{input: fmt.Sprintf("%d\n%s", t, b.String())})
	}

	return tests
}
