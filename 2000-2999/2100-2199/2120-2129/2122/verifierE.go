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

type test struct {
	input string
}

type caseData struct {
	n, k   int
	top    []int
	bottom []int
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

func formatCase(cs caseData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", cs.n, cs.k))
	for i, v := range cs.top {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range cs.bottom {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func formatInput(cases []caseData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, cs := range cases {
		sb.WriteString(formatCase(cs))
	}
	return sb.String()
}

func randomCase(rng *rand.Rand, nMax int) caseData {
	n := rng.Intn(nMax) + 1
	k := rng.Intn(500) + 1
	top := make([]int, n)
	bot := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(3) == 0 {
			top[i] = -1
		} else {
			top[i] = rng.Intn(k) + 1
		}
		if rng.Intn(3) == 0 {
			bot[i] = -1
		} else {
			bot[i] = rng.Intn(k) + 1
		}
	}
	return caseData{n: n, k: k, top: top, bottom: bot}
}

func generateTests() []test {
	var tests []test
	// Samples from statement
	sample := []caseData{
		{n: 4, k: 3, top: []int{2, 1, -1, 2}, bottom: []int{2, -1, 1, 3}},
		{n: 5, k: 4, top: []int{1, 3, -1, 4, 2}, bottom: []int{-1, 3, 4, 2, -1}},
		{n: 10, k: 10, top: []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, bottom: []int{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}},
	}
	tests = append(tests, test{input: formatInput(sample)})

	// Simple edge cases
	tests = append(tests, test{input: formatInput([]caseData{
		{n: 1, k: 1, top: []int{-1}, bottom: []int{-1}},
		{n: 1, k: 5, top: []int{3}, bottom: []int{-1}},
	})})

	rng := rand.New(rand.NewSource(2122))
	for i := 0; i < 40; i++ {
		tcCount := rng.Intn(3) + 1
		cases := make([]caseData, tcCount)
		sumN := 0
		for j := 0; j < tcCount; j++ {
			remain := 500 - sumN
			nMax := remain
			if nMax > 50 {
				nMax = 50
			}
			if nMax < 1 {
				nMax = 1
			}
			cases[j] = randomCase(rng, nMax)
			sumN += cases[j].n
		}
		tests = append(tests, test{input: formatInput(cases)})
	}

	// A larger stress case within limits
	stress := []caseData{
		randomCase(rand.New(rand.NewSource(7)), 120),
		randomCase(rand.New(rand.NewSource(8)), 120),
	}
	tests = append(tests, test{input: formatInput(stress)})

	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	candBin, candCleanup, err := prepareBinary(os.Args[1], "cand2122E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "2122E.go")
	refBin, refCleanup, err := prepareBinary(refPath, "ref2122E")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := generateTests()
	for i, tc := range tests {
		expected, err := runBinary(refBin, tc.input, 6*time.Second)
		if err != nil {
			fmt.Printf("Reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candBin, tc.input, 6*time.Second)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if expected != got {
			fmt.Printf("Wrong answer on test %d\nInput:\n%s\nExpected:\n%s\nGot:\n%s\n", i+1, tc.input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
