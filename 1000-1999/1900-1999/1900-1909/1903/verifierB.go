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

func getRefPath() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "1903B.go")
}

type testCase struct {
	n      int
	matrix [][]int
}

func generateInput() ([]byte, []testCase) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := 100
	var sb bytes.Buffer
	fmt.Fprintf(&sb, "%d\n", t)
	tests := make([]testCase, t)
	
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 2
		fmt.Fprintf(&sb, "%d\n", n)
		
		matrix := make([][]int, n)
		for r := 0; r < n; r++ {
			matrix[r] = make([]int, n)
		}
		
		if rng.Intn(2) == 0 {
			// YES case
			a := make([]int, n)
			for j := 0; j < n; j++ {
				a[j] = rng.Intn(1 << 30)
			}
			for r := 0; r < n; r++ {
				for c := 0; c < n; c++ {
					if r == c {
						matrix[r][c] = 0
					} else {
						matrix[r][c] = a[r] | a[c]
					}
				}
			}
		} else {
			// Random case (usually NO)
			for r := 0; r < n; r++ {
				for c := r + 1; c < n; c++ {
					val := rng.Intn(1 << 30)
					matrix[r][c] = val
					matrix[c][r] = val
				}
			}
		}
		
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				if c > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", matrix[r][c])
			}
			sb.WriteByte('\n')
		}
		tests[i] = testCase{n: n, matrix: matrix}
	}
	
	return sb.Bytes(), tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, tests := generateInput()

	refBin, cleanup, err := buildReference(getRefPath())
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	refAnswers, err := parseReference(refOut, tests)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := verifyCandidateOutput(candOut, tests, refAnswers); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "candidate output:")
		fmt.Fprintln(os.Stderr, candOut)
		os.Exit(1)
	}

	fmt.Println("Accepted")
}

func parseReference(out string, tests []testCase) ([]bool, error) {
	tokens := strings.Fields(out)
	idx := 0
	ans := make([]bool, len(tests))
	for ti, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("reference output ended early at test %d", ti+1)
		}
		word := strings.ToUpper(tokens[idx])
		idx++
		if word != "YES" && word != "NO" {
			return nil, fmt.Errorf("reference output test %d: expected YES/NO, got %q", ti+1, tokens[idx-1])
		}
		if word == "YES" {
			ans[ti] = true
			for j := 0; j < tc.n; j++ {
				if idx >= len(tokens) {
					return nil, fmt.Errorf("reference output missing values for test %d", ti+1)
				}
				if _, err := strconv.ParseInt(tokens[idx], 10, 64); err != nil {
					return nil, fmt.Errorf("reference output invalid integer at test %d: %v", ti+1, err)
				}
				idx++
			}
		}
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("reference output has unexpected extra tokens starting at %q", tokens[idx])
	}
	return ans, nil
}

func verifyCandidateOutput(out string, tests []testCase, ref []bool) error {
	tokens := strings.Fields(out)
	idx := 0
	for ti, tc := range tests {
		if idx >= len(tokens) {
			return fmt.Errorf("test %d: missing YES/NO token", ti+1)
		}
		word := strings.ToUpper(tokens[idx])
		idx++
		if word != "YES" && word != "NO" {
			return fmt.Errorf("test %d: expected YES/NO, got %q", ti+1, tokens[idx-1])
		}
		if word == "NO" {
			if ref[ti] {
				return fmt.Errorf("test %d: reported NO but reference found a solution", ti+1)
			}
			continue
		}
		if !ref[ti] {
			return fmt.Errorf("test %d: reported YES but reference claims impossible", ti+1)
		}
		if idx+tc.n > len(tokens) {
			return fmt.Errorf("test %d: expected %d numbers after YES, got %d", ti+1, tc.n, len(tokens)-idx)
		}
		arr := make([]int, tc.n)
		for j := 0; j < tc.n; j++ {
			val, err := strconv.Atoi(tokens[idx+j])
			if err != nil {
				return fmt.Errorf("test %d: invalid integer at position %d: %v", ti+1, j+1, err)
			}
			if val < 0 || val >= 1<<30 {
				return fmt.Errorf("test %d: value %d out of range [0, 2^30)", ti+1, val)
			}
			arr[j] = val
		}
		idx += tc.n
		if err := checkMatrix(tc.matrix, arr); err != nil {
			return fmt.Errorf("test %d: %v", ti+1, err)
		}
	}
	if idx != len(tokens) {
		return fmt.Errorf("unexpected extra tokens starting at %q", tokens[idx])
	}
	return nil
}

func checkMatrix(matrix [][]int, arr []int) error {
	n := len(arr)
	if len(matrix) != n {
		return fmt.Errorf("matrix size mismatch")
	}
	for i := 0; i < n; i++ {
		if len(matrix[i]) != n {
			return fmt.Errorf("matrix row %d has wrong length", i+1)
		}
		for j := i + 1; j < n; j++ {
			if (arr[i] | arr[j]) != matrix[i][j] {
				return fmt.Errorf("pair (%d,%d) expected %d but got %d", i+1, j+1, matrix[i][j], arr[i]|arr[j])
			}
		}
	}
	return nil
}

func buildReference(src string) (string, func(), error) {
	dir, err := os.MkdirTemp("", "verifier-1903B-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}