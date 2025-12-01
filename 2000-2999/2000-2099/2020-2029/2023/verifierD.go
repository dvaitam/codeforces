package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const referenceSource = "./2023D.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		refVal, err := parseFloat(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candVal, err := parseFloat(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}

		if !closeEnough(refVal, candVal) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %.10f got %.10f\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, refVal, candVal, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2023D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2023D.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
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

func buildTests() []testCase {
	var tests []testCase
	manual := []string{
		"3\n80 80\n70 100\n50 200\n",
		"2\n100 1\n100 1\n",
		"4\n1 100\n2 1000\n2 100\n3 1\n",
		"5\n34 80\n78 20\n99 19\n90 43\n65 79\n",
	}
	for idx, in := range manual {
		tests = append(tests, testCase{name: fmt.Sprintf("manual-%d", idx+1), input: in})
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, testCase{name: fmt.Sprintf("random-%d", i+1), input: randomTest(rng)})
	}
	return tests
}

func randomTest(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		p := rng.Intn(100) + 1
		maxW := 200000 / p
		if maxW < 1 {
			maxW = 1
		}
		w := rng.Intn(maxW) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, w))
	}
	return sb.String()
}

func parseFloat(out string) (float64, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseFloat(out, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q", out)
	}
	return val, nil
}

func closeEnough(expect, got float64) bool {
	diff := math.Abs(expect - got)
	den := math.Max(1.0, math.Abs(expect))
	return diff <= 1e-6+1e-6*den
}
