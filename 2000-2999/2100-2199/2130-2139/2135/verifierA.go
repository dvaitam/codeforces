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

type caseData struct {
	n   int
	arr []int
}

type test struct {
	input string
}

func prepareBinary(path, tag string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	bin := filepath.Join(os.TempDir(), tag)
	cmd := exec.Command("go", "build", "-o", bin, path)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", func() {}, fmt.Errorf("build %s: %v\n%s", path, err, out)
	}
	cleanup := func() { os.Remove(bin) }
	return bin, cleanup, nil
}

func runBinary(path, input string, timeout time.Duration) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	timer := time.AfterFunc(timeout, func() { cmd.Process.Kill() })
	err := cmd.Run()
	timer.Stop()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func formatInput(cases []caseData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomCase(rng *rand.Rand, n int) caseData {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n) + 1
	}
	return caseData{n: n, arr: arr}
}

func generateTests() []test {
	var tests []test

	// Deterministic small coverage
	tests = append(tests, test{input: formatInput([]caseData{
		{n: 1, arr: []int{1}},
		{n: 2, arr: []int{2, 2}},
		{n: 4, arr: []int{2, 2, 1, 1}},
		{n: 5, arr: []int{1, 3, 3, 3, 1}},
		{n: 5, arr: []int{4, 4, 4, 4, 2}},
	})})

	rng := rand.New(rand.NewSource(2135))
	for i := 0; i < 30; i++ {
		tcCnt := rng.Intn(4) + 1
		var cases []caseData
		sumN := 0
		for j := 0; j < tcCnt; j++ {
			n := rng.Intn(200) + 1
			sumN += n
			if sumN > 5000 {
				break
			}
			cases = append(cases, randomCase(rng, n))
		}
		if len(cases) > 0 {
			tests = append(tests, test{input: formatInput(cases)})
		}
	}

	// Stress near limits (sum n <= 2e5)
	largeN := 120000
	tests = append(tests, test{input: formatInput([]caseData{
		randomCase(rand.New(rand.NewSource(7)), largeN),
		randomCase(rand.New(rand.NewSource(8)), 40000),
	})})

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}

	candBin, candCleanup, err := prepareBinary(os.Args[1], "cand2135A")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "2135A.go")
	refBin, refCleanup, err := prepareBinary(refPath, "ref2135A")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runBinary(refBin, tc.input, 8*time.Second)
		if err != nil {
			fmt.Printf("Reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candBin, tc.input, 8*time.Second)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("Wrong answer on test %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
