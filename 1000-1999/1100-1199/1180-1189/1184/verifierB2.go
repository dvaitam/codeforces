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

func baseDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Dir(file)
}

type testCase struct {
	name  string
	input string
}

func runProgram(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("no output")
	}
	var val int64
	if _, err := fmt.Sscan(fields[0], &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer: %v", err)
	}
	return val, nil
}

func formatCase(n, m int, edges [][2]int, s, b int, k, h int64, ships []Ship, bases []Base) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", s, b, k, h))
	for _, sh := range ships {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", sh.x, sh.a, sh.f))
	}
	for _, bb := range bases {
		sb.WriteString(fmt.Sprintf("%d %d\n", bb.x, bb.d))
	}
	return sb.String()
}

type Ship struct {
	x int
	a int64
	f int64
}

type Base struct {
	x int
	d int64
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 60)
	// Simple manual case
	tests = append(tests, manualCase1())
	tests = append(tests, manualCase2())
	tests = append(tests, manualCase3())

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		tests = append(tests, randomCase(rng, i+1))
	}
	return tests
}

func manualCase1() testCase {
	n, m := 2, 1
	edges := [][2]int{{1, 2}}
	s, b := 1, 1
	k, h := int64(5), int64(7)
	ships := []Ship{{x: 1, a: 10, f: 1}}
	bases := []Base{{x: 2, d: 5}}
	input := formatCase(n, m, edges, s, b, k, h, ships, bases)
	return testCase{name: "manual_connected", input: input}
}

func manualCase2() testCase {
	n, m := 3, 2
	edges := [][2]int{{1, 2}, {2, 3}}
	s, b := 2, 1
	k, h := int64(3), int64(100)
	ships := []Ship{
		{x: 1, a: 5, f: 2},
		{x: 3, a: 1, f: 1},
	}
	bases := []Base{{x: 2, d: 3}}
	input := formatCase(n, m, edges, s, b, k, h, ships, bases)
	return testCase{name: "manual_unreachable", input: input}
}

func manualCase3() testCase {
	n, m := 4, 3
	edges := [][2]int{{1, 2}, {2, 3}, {3, 4}}
	s, b := 2, 2
	k, h := int64(8), int64(5)
	ships := []Ship{
		{x: 1, a: 10, f: 3},
		{x: 4, a: 4, f: 2},
	}
	bases := []Base{
		{x: 2, d: 5},
		{x: 3, d: 4},
	}
	input := formatCase(n, m, edges, s, b, k, h, ships, bases)
	return testCase{name: "manual_compare_costs", input: input}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	n := rng.Intn(20) + 1
	maxEdges := n * (n - 1) / 2
	m := 0
	if maxEdges > 0 {
		m = rng.Intn(maxEdges + 1)
	}
	allEdges := make([][2]int, 0, maxEdges)
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			allEdges = append(allEdges, [2]int{i, j})
		}
	}
	rng.Shuffle(len(allEdges), func(i, j int) {
		allEdges[i], allEdges[j] = allEdges[j], allEdges[i]
	})
	edges := allEdges[:m]

	s := rng.Intn(5) + 1
	if s > n {
		s = n
	}
	b := rng.Intn(5) + 1

	k := int64(rng.Intn(50) + 1)
	h := int64(rng.Intn(50) + 1)

	ships := make([]Ship, s)
	for i := 0; i < s; i++ {
		ships[i] = Ship{
			x: rng.Intn(n) + 1,
			a: int64(rng.Intn(50)),
			f: int64(rng.Intn(10)),
		}
	}
	bases := make([]Base, b)
	for i := 0; i < b; i++ {
		bases[i] = Base{
			x: rng.Intn(n) + 1,
			d: int64(rng.Intn(40)),
		}
	}
	input := formatCase(n, m, edges, s, b, k, h, ships, bases)
	return testCase{name: fmt.Sprintf("random_%d", idx), input: input}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	candidate, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve binary path: %v\n", err)
		os.Exit(1)
	}
	refPath := filepath.Join(baseDir(), "1184B2.go")

	tests := buildTests()
	for idx, tc := range tests {
		expOut, err := runProgram(refPath, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}
		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, expVal, gotVal, tc.input, expOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
