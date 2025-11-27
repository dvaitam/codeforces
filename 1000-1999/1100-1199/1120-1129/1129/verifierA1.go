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

const (
	refSourceA1 = "1000-1999/1100-1199/1120-1129/1129/1129A1.go"
	refBinaryA1 = "ref1129A1.bin"
	totalTests  = 80
)

type testCase struct {
	n   int
	m   int
	src []int
	dst []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath, cleanup, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		input := formatInput(tc)

		refOut, err := runProgram(refPath, input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		refVals, err := parseOutput(refOut, tc.n)
		if err != nil {
			fmt.Printf("failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			printInput(input)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", idx+1, err)
			printInput(input)
			os.Exit(1)
		}
		candVals, err := parseOutput(candOut, tc.n)
		if err != nil {
			fmt.Printf("failed to parse candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			printInput(input)
			os.Exit(1)
		}

		if len(refVals) != len(candVals) {
			fmt.Printf("test %d failed: expected %d numbers, got %d\n", idx+1, len(refVals), len(candVals))
			printInput(input)
			fmt.Println("Reference output:")
			fmt.Println(refOut)
			fmt.Println("Candidate output:")
			fmt.Println(candOut)
			os.Exit(1)
		}

		for i := range refVals {
			if refVals[i] != candVals[i] {
				fmt.Printf("test %d failed at position %d: expected %d, got %d\n", idx+1, i+1, refVals[i], candVals[i])
				printInput(input)
				fmt.Println("Reference output:")
				fmt.Println(refOut)
				fmt.Println("Candidate output:")
				fmt.Println(candOut)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref-1129A1-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1129A1.bin")
	srcDir, err := filepath.Abs(filepath.Dir(refSourceA1))
	if err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to resolve reference directory: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", bin, "./"+filepath.Base(refSourceA1))
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return bin, cleanup, nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(out string, n int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	res := make([]int64, n)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tc testCase) []byte {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.src[i], tc.dst[i]))
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 2, m: 1, src: []int{1}, dst: []int{2}},
		{n: 3, m: 2, src: []int{1, 2}, dst: []int{2, 3}},
		{n: 4, m: 3, src: []int{1, 2, 3}, dst: []int{2, 3, 4}},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests-2 {
		n := rng.Intn(50) + 2
		m := rng.Intn(100) + 1
		src := make([]int, m)
		dst := make([]int, m)
		for i := 0; i < m; i++ {
			src[i] = rng.Intn(n) + 1
			dst[i] = rng.Intn(n) + 1
			for dst[i] == src[i] {
				dst[i] = rng.Intn(n) + 1
			}
		}
		tests = append(tests, testCase{n: n, m: m, src: src, dst: dst})
	}
	tests = append(tests,
		testCase{
			n: 100,
			m: 200,
			src: func() []int {
				s := make([]int, 200)
				for i := 0; i < 200; i++ {
					s[i] = (i%100 + 1)
				}
				return s
			}(),
			dst: func() []int {
				d := make([]int, 200)
				for i := 0; i < 200; i++ {
					d[i] = ((i+1)%100 + 1)
				}
				return d
			}(),
		},
		testCase{
			n: 100,
			m: 200,
			src: func() []int {
				s := make([]int, 200)
				for i := 0; i < 200; i++ {
					s[i] = 100
				}
				return s
			}(),
			dst: func() []int {
				d := make([]int, 200)
				for i := 0; i < 200; i++ {
					d[i] = 1
				}
				return d
			}(),
		},
	)
	return tests
}

func printInput(in []byte) {
	fmt.Println("Input used:")
	fmt.Println(string(in))
}
