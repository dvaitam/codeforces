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
	n int
	k int64
	w int64
	h int64
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, err := buildReference()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n%s", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output parse error: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n%s", err, candOut)
		os.Exit(1)
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output parse error: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\nInput:\n%s\n", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1993F2-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1993F2.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
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
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	var tests []testCase
	sumN := 0
	add := func(tc testCase) {
		if sumN+tc.n > 1_000_000 {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	add(testCase{n: 2, k: 4, w: 2, h: 2, s: "UR"})
	add(testCase{n: 4, k: 2, w: 1, h: 1, s: "LLDD"})
	add(testCase{n: 6, k: 3, w: 3, h: 1, s: "RLRRRL"})
	add(testCase{n: 7, k: 123456789999, w: 3, h: 2, s: "ULULURD"})

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	moves := []byte{'L', 'R', 'U', 'D'}
	for len(tests) < 80 {
		n := rng.Intn(1000) + 1
		w := int64(rng.Intn(1_000_000) + 1)
		h := int64(rng.Intn(1_000_000) + 1)
		k := rng.Int63n(1_000_000_000_000) + 1
		sb := make([]byte, n)
		for i := 0; i < n; i++ {
			sb[i] = moves[rng.Intn(len(moves))]
		}
		add(testCase{n: n, k: k, w: w, h: h, s: string(sb)})
		if sumN >= 600_000 {
			break
		}
	}

	// Add one large test close to limits if space allows.
	if sumN < 1_000_000 {
		n := 1_000_000 - sumN
		if n < 1 {
			n = 1
		}
		if n > 100_000 {
			n = 100_000
		}
		sb := make([]byte, n)
		for i := 0; i < n; i++ {
			sb[i] = moves[i%4]
		}
		add(testCase{n: n, k: 1_000_000_000_000, w: 1_000_000, h: 1_000_000, s: string(sb)})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.k, tc.w, tc.h))
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int64, t)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q at position %d", tok, i+1)
		}
		res[i] = val
	}
	return res, nil
}
