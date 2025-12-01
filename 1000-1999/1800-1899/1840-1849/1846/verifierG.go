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
	refSource  = "./1846G.go"
	totalTests = 80
)

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		refAns, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		candAns, err := parseAnswers(candOut, tc.t)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}
		for i := 0; i < tc.t; i++ {
			if refAns[i] != candAns[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %d, got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refAns[i], candAns[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1846G-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1846G.bin")
	cmd := exec.Command("go", "build", "-o", bin, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("go build failed: %v\n%s", err, stderr.String())
	}
	return bin, func() { os.RemoveAll(dir) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseAnswers(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d", expected, len(fields))
	}
	ans := make([]int, expected)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		ans[i] = val
	}
	return ans, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "no_symptoms", input: "1\n2 0\n00\n", t: 1},
		{name: "simple_meds", input: "1\n2 2\n11\n1\n10\n01\n1\n01\n00\n", t: 1},
		{name: "impossible", input: "1\n2 1\n11\n1\n00\n11\n", t: 1},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests {
		tests = append(tests, randomCase(rng, len(tests)+1))
	}
	return tests
}

func randomCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(6) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		sb.WriteString(randomBits(rng, n))
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			d := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("%d\n", d))
			remove := randomBits(rng, n)
			add := randomBits(rng, n)
			for j := 0; j < n; j++ {
				if remove[j] == '1' && add[j] == '1' {
					add = add[:j] + "0" + add[j+1:]
				}
			}
			sb.WriteString(remove)
			sb.WriteByte('\n')
			sb.WriteString(add)
			sb.WriteByte('\n')
		}
	}
	return testCase{name: fmt.Sprintf("rand_%d", idx), input: sb.String(), t: t}
}

func randomBits(rng *rand.Rand, n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteByte('1')
		}
	}
	return sb.String()
}
