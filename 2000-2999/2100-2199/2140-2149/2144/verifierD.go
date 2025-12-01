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

const refSource = "./2144D.go"

type testCase struct {
	n   int
	y   int64
	arr []int
}

type testSuite struct {
	name  string
	cases []testCase
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
	dir, err := os.MkdirTemp("", "cf-2144D-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2144D.bin")
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

func parseOutput(raw string, expected int) ([]int64, error) {
	fields := strings.Fields(raw)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers, got %d", expected, len(fields))
	}
	res := make([]int64, expected)
	for i, tok := range fields {
		val, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", tok)
		}
		res[i] = val
	}
	return res, nil
}

func equalAnswers(a, b []int64) bool {
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.y))
		for i, v := range tc.arr {
			sb.WriteString(strconv.Itoa(v))
			if i+1 < len(tc.arr) {
				sb.WriteByte(' ')
			}
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildTests() []testSuite {
	sample := testSuite{
		name: "sample",
		cases: []testCase{
			{n: 5, y: 5, arr: []int{150, 150, 150, 50, 148}},
			{n: 3, y: 1000000000, arr: []int{42, 42, 42}},
			{n: 10, y: 5, arr: []int{11, 80, 88, 45, 1, 7, 3, 1, 9, 198}},
			{n: 4, y: 9999999, arr: []int{1, 1, 1, 1}},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var randomCases []testCase
	for i := 0; i < 10; i++ {
		n := rng.Intn(8) + 2
		y := int64(rng.Intn(50) + 1)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rng.Intn(200) + 1
		}
		randomCases = append(randomCases, testCase{n: n, y: y, arr: arr})
	}
	randomSuite := testSuite{name: "random", cases: randomCases}

	edges := testSuite{
		name: "edges",
		cases: []testCase{
			{n: 1, y: 1, arr: []int{1}},
			{n: 2, y: 1, arr: []int{200000, 200000}},
			{n: 2, y: 1_000_000_000, arr: []int{1, 2}},
		},
	}

	return []testSuite{sample, edges, randomSuite}
}
