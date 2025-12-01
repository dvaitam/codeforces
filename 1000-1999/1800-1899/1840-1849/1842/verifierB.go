package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	refSource  = "1000-1999/1800-1899/1840-1849/1842/1842B.go"
	totalTests = 80
)

type testCase struct {
	name  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
				fmt.Fprintf(os.Stderr, "test %d (%s) failed on case %d: expected %s, got %s\ninput:\n%sreference output:\n%s\ncandidate output:\n%s\n",
					idx+1, tc.name, i+1, refAns[i], candAns[i], tc.input, refOut, candOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "ref1842B-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1842B.bin")
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

func parseAnswers(out string, t int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(lines))
	}
	ans := make([]string, t)
	for i, s := range lines {
		ans[i] = strings.ToUpper(s)
		if ans[i] != "YES" && ans[i] != "NO" {
			return nil, fmt.Errorf("answer %d invalid: %q", i+1, s)
		}
	}
	return ans, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "trivial_zero", input: "1\n1 0\n0\n0\n0\n", t: 1},
		{name: "single_match", input: "1\n1 5\n5\n0\n0\n", t: 1},
		{name: "basic_example", input: "1\n4 7\n1 2 4 8\n5 1 3 0\n1 2 3 4\n", t: 1},
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
		n := rng.Intn(10) + 1
		if rng.Intn(5) == 0 {
			n = rng.Intn(50) + 1
		}
		x := rng.Intn(1_000_000_000 + 1)
		sb.WriteString(fmt.Sprintf("%d %d\n", n, x))
		for stack := 0; stack < 3; stack++ {
			for i := 0; i < n; i++ {
				if i > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprint(rng.Intn(1_000_000_001)))
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{name: fmt.Sprintf("rand_%d", idx), input: sb.String(), t: t}
}
