package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2111B.go"

var fib = []int{0, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144}

type testCase struct {
	name string
	n    int
	box  [][3]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate_binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	if err := validateOutputs(tests, refOut); err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}
	if err := validateOutputs(tests, candOut); err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2111B-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
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
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 64)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.box)))
		for _, b := range tc.box {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", b[0], b[1], b[2]))
		}
	}
	return sb.String()
}

func validateOutputs(tests []testCase, output string) error {
	lines := splitNonEmptyLines(output)
	if len(lines) != len(tests) {
		return fmt.Errorf("expected %d lines, got %d", len(tests), len(lines))
	}
	for i, line := range lines {
		res := strings.Join(strings.Fields(line), "")
		m := len(tests[i].box)
		if len(res) != m {
			return fmt.Errorf("case %d: expected string length %d, got %d", i+1, m, len(res))
		}
		for j, ch := range res {
			if ch != '0' && ch != '1' {
				return fmt.Errorf("case %d: invalid character %q", i+1, ch)
			}
			expect := canFit(tests[i].n, tests[i].box[j])
			if expect && ch == '0' {
				return fmt.Errorf("case %d: box %d marked 0 but cubes fit", i+1, j+1)
			}
			if !expect && ch == '1' {
				return fmt.Errorf("case %d: box %d marked 1 but cubes do not fit", i+1, j+1)
			}
		}
	}
	return nil
}

func splitNonEmptyLines(out string) []string {
	raw := strings.Split(out, "\n")
	lines := make([]string, 0, len(raw))
	for _, ln := range raw {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		lines = append(lines, ln)
	}
	return lines
}

// canFit returns true if all n Fibonacci cubes fit in a box with the given dimensions.
// The optimal arrangement tiles a 2D floor using the Fibonacci rectangle property:
// cubes f_1..f_n tile a f_n × f_{n+1} rectangle as a single layer.
// Necessary and sufficient condition (sorted a<=b<=c): a>=fib[n], b>=fib[n], c>=fib[n+1].
func canFit(n int, dims [3]int) bool {
	a, b, c := dims[0], dims[1], dims[2]
	if a > b {
		a, b = b, a
	}
	if b > c {
		b, c = c, b
	}
	if a > b {
		a, b = b, a
	}
	return a >= fib[n] && b >= fib[n] && c >= fib[n+1]
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "sample_like",
			n:    5,
			box: [][3]int{
				{3, 1, 2},
				{10, 10, 10},
				{9, 8, 13},
				{14, 7, 20},
			},
		},
		{
			name: "minimal_n",
			n:    2,
			box: [][3]int{
				{1, 1, 3},
				{2, 2, 2},
				{2, 3, 3},
			},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalM := 0
	for _, tc := range tests {
		totalM += len(tc.box)
	}

	for t := 0; t < 50 && totalM < 50000; t++ {
		n := rng.Intn(9) + 2 // 2..10
		m := rng.Intn(400) + 1
		boxes := make([][3]int, m)
		for i := 0; i < m; i++ {
			boxes[i] = [3]int{
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("rnd_small_%d", t+1),
			n:    n,
			box:  boxes,
		})
		totalM += m
	}

	if totalM < 200000 {
		// one heavy case
		m := 200000 - totalM
		if m > 5000 {
			m = 5000
		}
		boxes := make([][3]int, m)
		for i := 0; i < m; i++ {
			boxes[i] = [3]int{
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
				rng.Intn(150) + 1,
			}
		}
		tests = append(tests, testCase{
			name: "heavy_case",
			n:    10,
			box:  boxes,
		})
	}

	return tests
}
