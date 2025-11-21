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
	"strings"
	"time"
)

type testCase struct {
	n        int
	requests [][2]int
}

func parseInput(data string) ([]testCase, error) {
	reader := strings.NewReader(data)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]testCase, 0, t)
	for ; t > 0; t-- {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			return nil, fmt.Errorf("failed to read n m: %v", err)
		}
		req := make([][2]int, m)
		for i := 0; i < m; i++ {
			if _, err := fmt.Fscan(reader, &req[i][0], &req[i][1]); err != nil {
				return nil, fmt.Errorf("failed to read request: %v", err)
			}
		}
		cases = append(cases, testCase{n: n, requests: req})
	}
	return cases, nil
}

func runProgram(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func checkArrangement(vec string, tc testCase) error {
	if len(vec) != tc.n {
		return fmt.Errorf("expected string of length %d, got %d", tc.n, len(vec))
	}
	prefixR := make([]int, tc.n+1)
	prefixW := make([]int, tc.n+1)
	for i := 0; i < tc.n; i++ {
		prefixR[i+1] = prefixR[i]
		prefixW[i+1] = prefixW[i]
		if vec[i] == 'R' {
			prefixR[i+1]++
		} else if vec[i] == 'W' {
			prefixW[i+1]++
		} else {
			return fmt.Errorf("invalid character %q at position %d", vec[i], i+1)
		}
	}
	for _, req := range tc.requests {
		targetR, targetW := req[0], req[1]
		found := false
		for l := 0; l < tc.n && !found; l++ {
			for r := l + 1; r <= tc.n; r++ {
				rCount := prefixR[r] - prefixR[l]
				wCount := prefixW[r] - prefixW[l]
				if rCount == targetR && wCount == targetW {
					found = true
					break
				}
			}
		}
		if !found {
			return fmt.Errorf("no interval with %d R and %d W", targetR, targetW)
		}
	}
	return nil
}

func validateOutput(out string, tcs []testCase) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	for i, tc := range tcs {
		if !scanner.Scan() {
			return fmt.Errorf("missing output for test case %d", i+1)
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "IMPOSSIBLE" {
			return nil
		}
		if err := checkArrangement(line, tc); err != nil {
			return fmt.Errorf("case %d invalid: %v", i+1, err)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output lines detected")
	}
	return nil
}

func generateRandomTests() string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := 30
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		m := rng.Intn(5) + 1
		fmt.Fprintf(&sb, "%d %d\n", n, m)
		for j := 0; j < m; j++ {
			r := rng.Intn(n + 1)
			w := rng.Intn(n + 1 - r)
			if r == 0 && w == 0 {
				if rng.Intn(2) == 0 {
					r = 1
				} else {
					w = 1
				}
			}
			fmt.Fprintf(&sb, "%d %d\n", r, w)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierM.go /path/to/binary")
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1662M.go")

	input := generateRandomTests()
	tcs, err := parseInput(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse generated cases: %v\n", err)
		os.Exit(1)
	}

	refOut, err := runProgram(refPath, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validateOutput(refOut, tcs); err != nil {
		fmt.Fprintf(os.Stderr, "reference output invalid: %v\n", err)
		os.Exit(1)
	}

	out, err := runProgram(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := validateOutput(out, tcs); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\ninput:\n%soutput:\n%s", err, input, out)
		os.Exit(1)
	}
	fmt.Println("all tests passed")
}
