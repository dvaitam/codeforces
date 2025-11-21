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

const referenceSource = "1000-1999/1200-1299/1270-1279/1277/1277B.go"

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
	for len(tests) < 150 {
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

		t, err := parseTestCount(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error parsing test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}

		refVals, err := parseOutput(refOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOut)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		for i := 0; i < t; i++ {
			if refVals[i] != candVals[i] {
				fmt.Fprintf(os.Stderr, "test %d failed on case %d: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, i+1, refVals[i], candVals[i], input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceBinary() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1277B-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1277B.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, output.String())
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
		// sample from the statement
		"4\n6\n40 6 40 3 20 1\n1\n1024\n4\n2 4 8 16\n3\n3 1 7\n",
		buildSingleCase([]int{1}),
		buildSingleCase([]int{2}),
		buildSingleCase([]int{16, 8, 4, 2, 1}),
		buildSingleCase([]int{3, 5, 7}),
		"2\n5\n1 2 3 4 5\n5\n1024 512 256 128 64\n",
	}
}

func buildSingleCase(arr []int) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	writeArray(&sb, arr)
	return sb.String()
}

func writeArray(sb *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
}

func randomTestInput(rng *rand.Rand) string {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	used := 0
	for i := 0; i < t; i++ {
		remainCases := t - i
		maxAvail := 200000 - used - (remainCases - 1)
		limit := maxAvail
		if limit > 5000 {
			limit = 5000
		}
		n := rng.Intn(limit) + 1
		used += n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			val := randomValue(rng)
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomValue(rng *rand.Rand) int {
	switch rng.Intn(5) {
	case 0:
		// high power of two
		return 1 << uint(rng.Intn(30))
	case 1:
		// odd number
		return rng.Intn(1e9/2)*2 + 1
	default:
		return rng.Intn(1_000_000_000) + 1
	}
}

func parseTestCount(input string) (int, error) {
	reader := strings.NewReader(input)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func parseOutput(out string, t int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d integers, got %d", t, len(tokens))
	}
	res := make([]int, t)
	for i, tok := range tokens {
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative answer %d", val)
		}
		res[i] = val
	}
	return res, nil
}
