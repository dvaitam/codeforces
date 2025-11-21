package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const referenceSource = "1000-1999/1500-1599/1530-1539/1532/1532B.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := manualTests()
	for len(tests) < 200 {
		tests = append(tests, randomTestInput(rng))
	}

	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOut)
			os.Exit(1)
		}
		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}

		t, err := readTestCount(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}

		expect, err := parseOutput(refOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOut)
			os.Exit(1)
		}
		got, err := parseOutput(candOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		for i := 0; i < t; i++ {
			if expect[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d failed on query %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, i+1, expect[i], got[i], input, refOut, candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceBinary() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1532B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1532B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, buf.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func manualTests() []string {
	return []string{
		// sample from statement
		"6\n5 2 3\n100 1 4\n1 10 5\n1000000000 1 6\n1 1 1000000000\n1 1 999999999\n",
		buildSingleQuery(1, 1, 1),
		buildSingleQuery(5, 5, 5),
		buildSingleQuery(10, 1, 2),
		buildSingleQuery(10, 1, 1),
	}
}

func buildSingleQuery(a, b, k int64) string {
	return fmt.Sprintf("1\n%d %d %d\n", a, b, k)
}

func randomTestInput(rng *rand.Rand) string {
	t := rng.Intn(30) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		a := randomValue(rng)
		b := randomValue(rng)
		k := randomValue(rng)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, k))
	}
	return sb.String()
}

func randomValue(rng *rand.Rand) int64 {
	switch rng.Intn(5) {
	case 0:
		return int64(rng.Intn(1000) + 1)
	case 1:
		return int64(rng.Intn(1_000_000) + 1)
	default:
		return int64(rng.Intn(1_000_000_000) + 1)
	}
}

func readTestCount(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func parseOutput(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}
