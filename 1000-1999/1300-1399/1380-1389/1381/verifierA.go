package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type testCase struct {
	n int
	a string
	b string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(1381))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		var sb1, sb2 strings.Builder
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb1.WriteByte('0')
			} else {
				sb1.WriteByte('1')
			}
			if rng.Intn(2) == 0 {
				sb2.WriteByte('0')
			} else {
				sb2.WriteByte('1')
			}
		}
		tests[i] = testCase{n, sb1.String(), sb2.String()}
	}
	return tests
}

func apply(arr []byte, p int) {
	for i := 0; i < p/2; i++ {
		arr[i], arr[p-1-i] = arr[p-1-i], arr[i]
	}
	for i := 0; i < p; i++ {
		if arr[i] == '0' {
			arr[i] = '1'
		} else {
			arr[i] = '0'
		}
	}
}

func verify(n int, a, b string, ops []int) bool {
	s := []byte(a)
	for _, v := range ops {
		if v < 1 || v > n {
			return false
		}
		apply(s, v)
	}
	return string(s) == b
}

func parseOutput(out string, n int) ([]int, error) {
	r := strings.NewReader(out)
	var k int
	if _, err := fmt.Fscan(r, &k); err != nil {
		return nil, err
	}
	if k < 0 || k > 3*n {
		return nil, fmt.Errorf("invalid k %d", k)
	}
	ops := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(r, &ops[i]); err != nil {
			return nil, err
		}
	}
	return ops, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := genTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", tc.n, tc.a, tc.b)
		out, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		ops, err := parseOutput(out, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error: %v\noutput:%s\n", i+1, err, out)
			os.Exit(1)
		}
		if !verify(tc.n, tc.a, tc.b, ops) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:%soutput:%s\n", i+1, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
