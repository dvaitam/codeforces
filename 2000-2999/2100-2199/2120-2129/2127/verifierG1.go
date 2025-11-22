package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type caseData struct {
	n    int
	perm []int
}

type test struct {
	input string
	cases []caseData
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
		for i, v := range cs.perm {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func derangement(n int, rng *rand.Rand) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) { p[i], p[j] = p[j], p[i] })
	for i := 0; i < n; i++ {
		if p[i] == i+1 {
			j := (i + 1) % n
			p[i], p[j] = p[j], p[i]
		}
	}
	return p
}

func generateTests() []test {
	var tests []test

	// Simple fixed cases
	tests = append(tests, test{cases: []caseData{
		{n: 4, perm: []int{2, 1, 4, 3}},
	}})
	tests = append(tests, test{cases: []caseData{
		{n: 5, perm: []int{2, 3, 1, 5, 4}},
		{n: 6, perm: []int{3, 1, 4, 6, 2, 5}},
	}})

	rng := rand.New(rand.NewSource(2127))
	totalN2 := 0
	for len(tests) < 40 && totalN2 < 9000 {
		tcCount := rng.Intn(3) + 1
		var cases []caseData
		for j := 0; j < tcCount; j++ {
			remain := 10000 - totalN2
			if remain < 16 { // 4^2 minimal
				break
			}
			maxN := int(math.Sqrt(float64(remain)))
			if maxN > 100 {
				maxN = 100
			}
			n := rng.Intn(maxN-3) + 4 // ensure at least 4
			perm := derangement(n, rng)
			cases = append(cases, caseData{n: n, perm: perm})
			totalN2 += n * n
		}
		if len(cases) > 0 {
			tests = append(tests, test{cases: cases})
		}
	}

	for i := range tests {
		tests[i].input = formatInput(tests[i].cases)
	}
	return tests
}

func compareOutputs(cases []caseData, out string) error {
	tokens := strings.Fields(out)
	pos := 0
	for idx, cs := range cases {
		if pos+cs.n > len(tokens) {
			return fmt.Errorf("case %d: not enough numbers in output", idx+1)
		}
		for i := 0; i < cs.n; i++ {
			v := atoi(tokens[pos+i])
			if v != cs.perm[i] {
				return fmt.Errorf("case %d: position %d expected %d got %d", idx+1, i+1, cs.perm[i], v)
			}
		}
		pos += cs.n
	}
	if pos != len(tokens) {
		return fmt.Errorf("extra tokens in output")
	}
	return nil
}

func atoi(s string) int {
	sign := 1
	i := 0
	if len(s) > 0 && s[0] == '-' {
		sign = -1
		i++
	}
	val := 0
	for ; i < len(s); i++ {
		val = val*10 + int(s[i]-'0')
	}
	return sign * val
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/binary")
		os.Exit(1)
	}

	candBin, candCleanup, err := prepareBinary(os.Args[1], "cand2127G1")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer candCleanup()

	_, file, _, _ := runtime.Caller(0)
	refPath := filepath.Join(filepath.Dir(file), "2127G1.go")
	refBin, refCleanup, err := prepareBinary(refPath, "ref2127G1")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer refCleanup()

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runBinary(refBin, tc.input, 4*time.Second)
		if err != nil {
			fmt.Printf("Reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(tc.cases, refOut); err != nil {
			fmt.Printf("Reference output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		got, err := runBinary(candBin, tc.input, 4*time.Second)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := compareOutputs(tc.cases, got); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
