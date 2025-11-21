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
	n      int
	top    []int
	bottom []int
}

func buildReference() (string, error) {
	path := "./2163C_ref.bin"
	cmd := exec.Command("go", "build", "-o", path, "2163C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return path, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func writeRow(sb *strings.Builder, vals []int) {
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
}

func casesToInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		writeRow(&sb, tc.top)
		writeRow(&sb, tc.bottom)
	}
	return sb.String()
}

func permutationCase(rng *rand.Rand, n int) testCase {
	vals := make([]int, 2*n)
	for i := range vals {
		vals[i] = i + 1
	}
	rng.Shuffle(len(vals), func(i, j int) {
		vals[i], vals[j] = vals[j], vals[i]
	})
	top := append([]int(nil), vals[:n]...)
	bottom := append([]int(nil), vals[n:]...)
	return testCase{n: n, top: top, bottom: bottom}
}

func randomMultiCase(rng *rand.Rand, n int) testCase {
	top := make([]int, n)
	bottom := make([]int, n)
	for i := 0; i < n; i++ {
		top[i] = rng.Intn(2*n) + 1
		bottom[i] = rng.Intn(2*n) + 1
	}
	return testCase{n: n, top: top, bottom: bottom}
}

func structuredCase(n int) testCase {
	top := make([]int, n)
	bottom := make([]int, n)
	for i := 0; i < n; i++ {
		top[i] = i + 1
		bottom[i] = 2*n - i
	}
	return testCase{n: n, top: top, bottom: bottom}
}

func layeredBatch(total int, rng *rand.Rand) string {
	var cases []testCase
	remaining := total
	for remaining > 0 {
		maxN := 100
		if remaining < maxN {
			maxN = remaining
		}
		var n int
		if remaining <= 2 {
			n = remaining
		} else {
			n = rng.Intn(maxN-1) + 2
			if n > remaining {
				n = remaining
			}
			if remaining-n == 1 {
				n++
			}
		}
		if n < 2 {
			n = 2
		}
		if rng.Intn(2) == 0 {
			cases = append(cases, permutationCase(rng, n))
		} else {
			cases = append(cases, randomMultiCase(rng, n))
		}
		remaining -= n
	}
	return casesToInput(cases)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var refPath string
	fail := func(format string, args ...interface{}) {
		if refPath != "" {
			_ = os.Remove(refPath)
		}
		fmt.Fprintf(os.Stderr, format+"\n", args...)
		os.Exit(1)
	}

	var err error
	refPath, err = buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer os.Remove(refPath)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []string{}

	tests = append(tests, casesToInput([]testCase{{
		n:      2,
		top:    []int{1, 2},
		bottom: []int{3, 4},
	}}))

	tests = append(tests, casesToInput([]testCase{{
		n:      3,
		top:    []int{1, 1, 1},
		bottom: []int{6, 6, 6},
	}, {
		n:      5,
		top:    []int{1, 5, 2, 6, 3},
		bottom: []int{4, 7, 8, 9, 10},
	}}))

	tests = append(tests, casesToInput([]testCase{structuredCase(10)}))

	for i := 0; i < 200; i++ {
		t := rng.Intn(5) + 1
		cases := make([]testCase, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(9) + 2
			if rng.Intn(2) == 0 {
				cases[j] = permutationCase(rng, n)
			} else {
				cases[j] = randomMultiCase(rng, n)
			}
		}
		tests = append(tests, casesToInput(cases))
	}

	for i := 0; i < 50; i++ {
		t := rng.Intn(3) + 1
		cases := make([]testCase, t)
		for j := 0; j < t; j++ {
			n := rng.Intn(80) + 20
			if rng.Intn(2) == 0 {
				cases[j] = permutationCase(rng, n)
			} else {
				cases[j] = randomMultiCase(rng, n)
			}
		}
		tests = append(tests, casesToInput(cases))
	}

	tests = append(tests, layeredBatch(200000, rng))
	tests = append(tests, casesToInput([]testCase{structuredCase(200000)}))

	for idx, input := range tests {
		expect, err := runProgram(refPath, input)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		got, err := runProgram(bin, input)
		if err != nil {
			fail("test %d: runtime error: %v\ninput:\n%s", idx+1, err, input)
		}
		if expect != got {
			fail("test %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s", idx+1, input, expect, got)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
