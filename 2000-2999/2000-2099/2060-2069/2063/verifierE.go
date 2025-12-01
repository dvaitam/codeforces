package main

import (
	"bufio"
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

const refSource = "./2063E.go"

type testCase struct {
	n   int
	edg [][2]int
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2063E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1},
		{n: 2, edg: [][2]int{{1, 2}}},
		{n: 3, edg: [][2]int{{1, 2}, {1, 3}}},
		{n: 3, edg: [][2]int{{1, 2}, {2, 3}}},
		{n: 5, edg: [][2]int{{1, 2}, {1, 3}, {2, 4}, {2, 5}}},
		{n: 6, edg: [][2]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}, {1, 6}}},
		{n: 7, edg: [][2]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}}},
	}
}

func randomTree(rng *rand.Rand, n int) testCase {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		parent := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{parent, v})
	}
	return testCase{n: n, edg: edges}
}

func heavyChain(n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return testCase{n: n, edg: edges}
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, e := range tc.edg {
			sb.WriteString(strconv.Itoa(e[0]))
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(e[1]))
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func parseOutput(out string, t int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]string, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		res[i] = sc.Text()
	}
	if sc.Scan() {
		return nil, fmt.Errorf("extra output detected after %d testcases", t)
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		tests = append(tests, randomTree(rng, rng.Intn(15)+2))
	}
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTree(rng, rng.Intn(150)+10))
	}
	for i := 0; i < 20; i++ {
		tests = append(tests, randomTree(rng, rng.Intn(2000)+200))
	}
	tests = append(tests, heavyChain(200000))
	tests = append(tests, randomTree(rng, 100000))

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseOutput(wantOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if want[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\nn=%d edges=%v\n",
				i+1, want[i], got[i], tests[i].n, tests[i].edg)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
