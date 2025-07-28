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
)

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag)
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n int, s string) []int {
	zeros := make([]int, 0, n)
	for i, ch := range s {
		if ch == '0' {
			zeros = append(zeros, i+1) // positions 1-indexed
		}
	}
	k := len(zeros)
	prefix := make([]int, k+1)
	for i := 0; i < k; i++ {
		prefix[i+1] = prefix[i] + zeros[i]
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		if i > k {
			res[i-1] = -1
			continue
		}
		targetSum := (n-i+1)*i + i*(i-1)/2
		sumZeros := prefix[k] - prefix[k-i]
		cost := targetSum - sumZeros
		res[i-1] = cost
	}
	return res
}

func generateTests() []struct {
	n int
	s string
} {
	rng := rand.New(rand.NewSource(2))
	tests := make([]struct {
		n int
		s string
	}, 0, 100)
	tests = append(tests, struct {
		n int
		s string
	}{1, "0"}, struct {
		n int
		s string
	}{1, "1"})
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		b := make([]byte, n)
		for i := range b {
			if rng.Intn(2) == 0 {
				b[i] = '0'
			} else {
				b[i] = '1'
			}
		}
		tests = append(tests, struct {
			n int
			s string
		}{n, string(b)})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candB")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n", tc.n, tc.s)
		expInts := solveCase(tc.n, tc.s)
		expStrs := make([]string, len(expInts))
		for j, v := range expInts {
			expStrs[j] = strconv.Itoa(v)
		}
		expected := strings.Join(expStrs, " ")
		output, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("case %d failed:\ninput:\n%d\n%s\nexpected: %s\ngot: %s\n", i+1, tc.n, tc.s, expected, strings.TrimSpace(output))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
