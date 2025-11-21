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

const referenceSource = "1000-1999/1600-1699/1660-1669/1662/1662N.go"

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierN.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := manualTests()
	for len(tests) < 200 {
		tests = append(tests, randomTest(rng))
	}

	for idx, input := range tests {
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOut)
			os.Exit(1)
		}
		expect, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		got, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d: %v\ninput:\n%soutput:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		if expect != got {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
				idx+1, expect, got, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-1662N-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref1662N.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
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
		// n = 2 impossible
		"2\n1 2\n3 4\n",
		// n = 2 valid
		"2\n1 3\n4 2\n",
		// n = 3 sample-like
		"3\n1 2 3\n4 5 6\n7 8 9\n",
	}
}

func randomTest(rng *rand.Rand) string {
	n := rng.Intn(4) + 2 // 2..5
	total := n * n
	values := make([]int, total)
	for i := 0; i < total; i++ {
		values[i] = i + 1
	}
	for i := total - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		values[i], values[j] = values[j], values[i]
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	idx := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(values[idx]))
			idx++
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(tokens))
	}
	val, err := strconv.ParseInt(tokens[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", tokens[0])
	}
	if val < 0 {
		return 0, fmt.Errorf("negative answer %d", val)
	}
	return val, nil
}
