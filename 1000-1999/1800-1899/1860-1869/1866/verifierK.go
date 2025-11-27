package main

import (
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

func runProgram(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return out.String(), nil
}

func parseOutput(out string, q int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) < q {
		return nil, fmt.Errorf("expected %d outputs, got %d", q, len(fields))
	}
	ans := make([]int64, q)
	for i := 0; i < q; i++ {
		if _, err := fmt.Sscan(fields[i], &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to parse integer at query %d: %v", i+1, err)
		}
	}
	return ans, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}

	_, src, _, _ := runtime.Caller(0)
	baseDir := filepath.Dir(src)
	refPath := filepath.Join(baseDir, "1866K.go")

	// If stdin is provided, use it as a single test; otherwise generate a few random tests.
	stdinData := readAllStdin()
	var tests [][]byte
	if len(stdinData) > 0 {
		tests = append(tests, stdinData)
	} else {
		tests = append(tests, buildSample())
		tests = append(tests, buildChain(5))
		tests = append(tests, buildRandom(8, 6, 1))
		tests = append(tests, buildRandom(12, 10, 2))
		tests = append(tests, buildRandom(20, 15, 3))
		tests = append(tests, buildRandom(30, 20, 4))
		tests = append(tests, buildRandom(50, 30, 5))
		// bulk randoms to reach at least 100 cases
		seed := time.Now().UnixNano()
		for i := 0; i < 93; i++ {
			tests = append(tests, buildRandom(15+i%10, 10+i%5, seed+int64(i)))
		}
	}

	for ti, inputData := range tests {
		reader := strings.NewReader(string(inputData))
		var n int
		if _, err := fmt.Fscan(reader, &n); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read n on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}
		for i := 0; i < n-1; i++ {
			var u, v int
			var w int64
			fmt.Fscan(reader, &u, &v, &w)
		}
		var q int
		if _, err := fmt.Fscan(reader, &q); err != nil {
			fmt.Fprintf(os.Stderr, "failed to read q on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}

		refOut, err := runProgram(refPath, inputData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}
		expected, err := parseOutput(refOut, q)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}

		out, err := runProgram(target, inputData)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}
		ans, err := parseOutput(out, q)
	 if err != nil {
			fmt.Fprintf(os.Stderr, "target output parse error on test %d: %v\n", ti+1, err)
			os.Exit(1)
		}

		for i := 0; i < q; i++ {
			if ans[i] != expected[i] {
				fmt.Fprintf(os.Stderr, "wrong answer on test %d query %d: expected %d got %d\n", ti+1, i+1, expected[i], ans[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("all tests passed")
}

func readAllStdin() []byte {
	data, _ := os.ReadFile("/dev/stdin")
	return data
}

// Simple deterministic samples.
func buildSample() []byte {
	var sb strings.Builder
	// small star
	sb.WriteString("5\n")
	sb.WriteString("1 2 3\n")
	sb.WriteString("1 3 4\n")
	sb.WriteString("2 4 2\n")
	sb.WriteString("3 5 1\n")
	sb.WriteString("3\n")
	sb.WriteString("1 1\n")
	sb.WriteString("2 3\n")
	sb.WriteString("4 2\n")
	return []byte(sb.String())
}

func buildChain(n int) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", i, i+1, i))
	}
	sb.WriteString("3\n")
	sb.WriteString("1 1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, 1))
	sb.WriteString(fmt.Sprintf("%d %d\n", n/2, 2))
	return []byte(sb.String())
}

func buildRandom(n, q int, seed int64) []byte {
	rand.Seed(seed)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		w := rand.Intn(9) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", p, i, w))
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		v := rand.Intn(n) + 1
		k := rand.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", v, k))
	}
	return []byte(sb.String())
}
