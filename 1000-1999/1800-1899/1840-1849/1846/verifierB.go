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
	refSource  = "1846B.go"
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
	dir, err := os.MkdirTemp("", "ref1846B-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(dir, "ref1846B.bin")
	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", bin, source)
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
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(lines))
	}
	ans := make([]string, t)
	for i, s := range lines {
		s = strings.ToUpper(s)
		if s != "X" && s != "O" && s != "+" && s != "DRAW" {
			return nil, fmt.Errorf("invalid answer %q", s)
		}
		ans[i] = s
	}
	return ans, nil
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "simple_x", input: "1\nXXX\n...\n...\n", t: 1},
		{name: "simple_o", input: "1\nOOO\n...\n...\n", t: 1},
		{name: "simple_plus", input: "1\n+++\n...\n...\n", t: 1},
		{name: "draw_empty", input: "1\n...\n...\n...\n", t: 1},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < totalTests {
		tests = append(tests, randomCase(rng, len(tests)+1))
	}
	return tests
}

func randomCase(rng *rand.Rand, idx int) testCase {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for c := 0; c < t; c++ {
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				v := rng.Intn(4)
				switch v {
				case 0:
					sb.WriteByte('X')
				case 1:
					sb.WriteByte('O')
				case 2:
					sb.WriteByte('+')
				default:
					sb.WriteByte('.')
				}
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{name: fmt.Sprintf("rand_%d", idx), input: sb.String(), t: t}
}
