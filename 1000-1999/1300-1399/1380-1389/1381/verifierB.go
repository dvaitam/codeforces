package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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

func buildOracle() (string, error) {
	tmp, err := os.CreateTemp("", "oracle*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if out, err := exec.Command("go", "build", "-o", tmp.Name(), "1381B.go").CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return tmp.Name(), nil
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
	p []int
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(10) + 1
		perm := rng.Perm(2 * n)
		for j := 0; j < 2*n; j++ {
			perm[j]++
		}
		tests[i] = testCase{n, perm}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for j, v := range tc.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle fail case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(want) {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
