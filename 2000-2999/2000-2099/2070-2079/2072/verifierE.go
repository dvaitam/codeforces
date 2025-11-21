package main

import (
	"bufio"
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

const (
	maxPoints  = 500
	coordLimit = 1_000_000_000
)

type testCase struct {
	k    int
	name string
}

type point struct {
	x int64
	y int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2072E-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleE")
	cmd := exec.Command("go", "build", "-o", outPath, "2072E.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
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

func generateTests() []testCase {
	base := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		15, 21, 28, 36, 45,
		55, 66, 78, 91, 100,
		500, 1000, 5000, 9999,
		12345, 22222, 32123, 44444, 54321,
		75000, 90000, 99999, 100000,
	}
	tests := make([]testCase, 0, len(base)+120)
	for idx, v := range base {
		tests = append(tests, testCase{
			k:    v,
			name: fmt.Sprintf("base_%d", idx+1),
		})
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 120; i++ {
		tests = append(tests, testCase{
			k:    rng.Intn(100000 + 1),
			name: fmt.Sprintf("random_%d", i+1),
		})
	}
	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.k))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func validateOutput(output string, tests []testCase) error {
	reader := bufio.NewReader(strings.NewReader(output))
	for idx, tc := range tests {
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return fmt.Errorf("test %d (%s): failed to read n: %v", idx+1, tc.name, err)
		}
		if n < 0 || n > maxPoints {
			return fmt.Errorf("test %d (%s): n=%d out of range [0,%d]", idx+1, tc.name, n, maxPoints)
		}
		points := make([]point, n)
		seen := make(map[point]struct{}, n)
		for i := 0; i < n; i++ {
			var x, y int64
			if _, err := fmt.Fscan(reader, &x, &y); err != nil {
				return fmt.Errorf("test %d (%s): failed to read point %d: %v", idx+1, tc.name, i+1, err)
			}
			if x < -coordLimit || x > coordLimit || y < -coordLimit || y > coordLimit {
				return fmt.Errorf("test %d (%s): point %d=(%d,%d) exceeds coordinate limits ±%d", idx+1, tc.name, i+1, x, y, coordLimit)
			}
			p := point{x: x, y: y}
			if _, exists := seen[p]; exists {
				return fmt.Errorf("test %d (%s): duplicate point (%d,%d)", idx+1, tc.name, x, y)
			}
			seen[p] = struct{}{}
			points[i] = p
		}

		pairs := 0
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if points[i].x == points[j].x || points[i].y == points[j].y {
					pairs++
				}
			}
		}
		if pairs != tc.k {
			return fmt.Errorf("test %d (%s): expected %d aligned pairs, got %d", idx+1, tc.name, tc.k, pairs)
		}
	}

	// Ensure no extra tokens remain (ignoring whitespace).
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("extra output detected after all test cases (e.g. %q)", extra)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	input := buildInput(tests)

	if _, err := runBinary(oracle, input); err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed on generated tests: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	output, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	if err := validateOutput(output, tests); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%soutput:\n%s", err, input, output)
		os.Exit(1)
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
