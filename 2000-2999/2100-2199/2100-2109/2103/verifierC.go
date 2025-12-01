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

const refSource = "./2103C.go"

type testCase struct {
	n int
	k int
	a []int
}

type testSuite struct {
	name  string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
	dir, err := os.MkdirTemp("", "cf-2103C-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2103C.bin")
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
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(fields))
	}
	res := make([]bool, expected)
	for i, tok := range fields {
		upper := strings.ToUpper(tok)
		switch upper {
		case "YES":
			res[i] = true
		case "NO":
			res[i] = false
		default:
			return nil, fmt.Errorf("invalid token %q (only YES/NO allowed)", tok)
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
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.k))
		for i, v := range c.a {
			sb.WriteString(strconv.Itoa(v))
			if i+1 < len(c.a) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testSuite {
	samples := []testCase{
		{n: 3, k: 2, a: []int{3, 2, 1}},
		{n: 3, k: 1, a: []int{3, 2, 1}},
		{n: 3, k: 13, a: []int{2, 1, 6}},
		{n: 8, k: 7, a: []int{10, 7, 12, 16, 3, 15, 6, 11}},
		{n: 6, k: 8, a: []int{7, 11, 12, 4, 9, 17}},
		{n: 3, k: 500000000, a: []int{1000, 1000000000, 1000}},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomCases []testCase
	for i := 0; i < 6; i++ {
		randomCases = append(randomCases, randomCase(rng, 5+i*3))
	}

	mixed := testSuite{name: "samples", cases: samples}
	randomSuite := testSuite{name: "random", cases: randomCases}

	return []testSuite{mixed, randomSuite}
}

func randomCase(rng *rand.Rand, n int) testCase {
	if n < 3 {
		n = 3
	}
	k := rng.Intn(20) + 1
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(30) + 1
	}
	return testCase{n: n, k: k, a: a}
}
