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

// refSource points to the local reference solution to avoid GOPATH resolution.
const refSource = "2131E.go"

type testCase struct {
	n int
	a []int
	b []int
}

type testSuite struct {
	name  string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	suites := buildTests()

	for idx, suite := range suites {
		input := buildInput(suite.cases)

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, refOut)
			os.Exit(1)
		}
		expected, err := parseOutput(refOut, len(suite.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, candOut)
			os.Exit(1)
		}
		got, err := parseOutput(candOut, len(suite.cases))
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, suite.name, err, input, candOut)
			os.Exit(1)
		}
		if !equalAnswers(expected, got) {
			fmt.Fprintf(os.Stderr, "candidate mismatch on test %d (%s)\ninput:\n%soutput:\n%s", idx+1, suite.name, input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d test suites passed\n", len(suites))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2131E-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2131E.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(raw string, expected int) ([]bool, error) {
	fields := strings.Fields(raw)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d YES/NO tokens, got %d", expected, len(fields))
	}
	res := make([]bool, expected)
	for i, tok := range fields {
		switch strings.ToUpper(tok) {
		case "YES":
			res[i] = true
		case "NO":
			res[i] = false
		default:
			return nil, fmt.Errorf("invalid token %q", tok)
		}
	}
	return res, nil
}

func equalAnswers(a, b []bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			sb.WriteString(strconv.Itoa(v))
			if i+1 < len(tc.a) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
		for i, v := range tc.b {
			sb.WriteString(strconv.Itoa(v))
			if i+1 < len(tc.b) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testSuite {
	sample := testSuite{
		name: "samples",
		cases: []testCase{
			{
				n: 5,
				a: []int{1, 2, 3, 4, 5},
				b: []int{3, 2, 7, 1, 5},
			},
			{
				n: 5,
				a: []int{0, 0, 1, 0, 1},
				b: []int{0, 0, 1, 0, 0},
			},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomCases []testCase
	for i := 0; i < 8; i++ {
		n := rng.Intn(8) + 2
		a := make([]int, n)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(16)
			b[j] = rng.Intn(16)
		}
		randomCases = append(randomCases, testCase{n: n, a: a, b: b})
	}
	randomSuite := testSuite{name: "random", cases: randomCases}

	return []testSuite{sample, randomSuite}
}
