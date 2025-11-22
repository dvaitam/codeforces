package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// computeMax replicates the reference greedy to obtain the optimal value.
func computeMax(a, b []int) int {
	n := len(a)
	lim := 2*n + 3
	seenA := make([]bool, lim)
	seenB := make([]bool, lim)

	for i := 0; i < n; i++ {
		x, y := a[i], b[i]
		gainOrig := 0
		if !seenA[x] {
			gainOrig++
		}
		if !seenB[y] {
			gainOrig++
		}

		gainSwap := 0
		if !seenA[y] {
			gainSwap++
		}
		if !seenB[x] {
			gainSwap++
		}

		if gainSwap > gainOrig {
			x, y = y, x
		}

		if x < lim {
			seenA[x] = true
		}
		if y < lim {
			seenB[y] = true
		}
	}

	ans := 0
	for i := 0; i < lim; i++ {
		if seenA[i] {
			ans++
		}
		if seenB[i] {
			ans++
		}
	}
	return ans
}

type testCase struct {
	raw string
}

func distinctCount(arr []int) int {
	m := make(map[int]struct{}, len(arr))
	for _, v := range arr {
		m[v] = struct{}{}
	}
	return len(m)
}

func parseAndJudge(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("failed to read t from input")
	}

	type task struct {
		n int
		a []int
		b []int
	}
	tasks := make([]task, t)
	for idx := 0; idx < t; idx++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return fmt.Errorf("input case %d: failed to read n", idx+1)
		}
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		tasks[idx] = task{n: n, a: a, b: b}
	}

	out := bufio.NewReader(strings.NewReader(output))
	for idx, tc := range tasks {
		var declared int
		if _, err := fmt.Fscan(out, &declared); err != nil {
			return fmt.Errorf("output case %d: missing declared answer", idx+1)
		}
		aOut := make([]int, tc.n)
		bOut := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(out, &aOut[i]); err != nil {
				return fmt.Errorf("output case %d: missing a[%d]", idx+1, i+1)
			}
		}
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(out, &bOut[i]); err != nil {
				return fmt.Errorf("output case %d: missing b[%d]", idx+1, i+1)
			}
		}

		for i := 0; i < tc.n; i++ {
			origA, origB := tc.a[i], tc.b[i]
			if !((aOut[i] == origA && bOut[i] == origB) || (aOut[i] == origB && bOut[i] == origA)) {
				return fmt.Errorf("output case %d: position %d invalid (a=%d,b=%d not from original pair %d,%d)", idx+1, i+1, aOut[i], bOut[i], origA, origB)
			}
			if aOut[i] < 1 || aOut[i] > 2*tc.n || bOut[i] < 1 || bOut[i] > 2*tc.n {
				return fmt.Errorf("output case %d: value out of allowed range at position %d", idx+1, i+1)
			}
		}

		val := distinctCount(aOut) + distinctCount(bOut)
		if val != declared {
			return fmt.Errorf("output case %d: declared value %d mismatches computed %d", idx+1, declared, val)
		}

		expected := computeMax(tc.a, tc.b)
		if val != expected {
			return fmt.Errorf("output case %d: value %d not optimal (expected %d)", idx+1, val, expected)
		}
	}

	var extra string
	if _, err := fmt.Fscan(out, &extra); err == nil {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2113))
	var tests []testCase

	// Sample from statement.
	sample := "3\n5\n1 2 4 4 4\n1 3 3 5 2\n7\n2 2 4 4 5 5 5\n1 3 3 2 1 6 6\n7\n12 3 3 4 5 6 4\n1 2 13 8 10 13 7\n"
	tests = append(tests, testCase{raw: sample})

	for len(tests) < 120 {
		t := rng.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for i := 0; i < t; i++ {
			n := rng.Intn(10) + 1
			fmt.Fprintf(&sb, "%d\n", n)
			for j := 0; j < n; j++ {
				fmt.Fprintf(&sb, "%d ", rng.Intn(2*n)+1)
			}
			sb.WriteByte('\n')
			for j := 0; j < n; j++ {
				fmt.Fprintf(&sb, "%d ", rng.Intn(2*n)+1)
			}
			sb.WriteByte('\n')
		}
		tests = append(tests, testCase{raw: sb.String()})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		got, err := runBinary(bin, tc.raw)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := parseAndJudge(tc.raw, got); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, tc.raw, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
