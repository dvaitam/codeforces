package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCase struct {
	desc  string
	input string
	t     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]
	ref, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Printf("Reference runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, refOut)
			os.Exit(1)
		}
		exp, err := parseAnswers(refOut, tc.t)
		if err != nil {
			fmt.Printf("Failed to parse reference output on test %d (%s): %v\nOutput:\n%s\n", i+1, tc.desc, err, refOut)
			os.Exit(1)
		}

		out, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		got, err := parseAnswers(out, tc.t)
		if err != nil {
			fmt.Printf("Failed to parse output on test %d (%s): %v\nInput:\n%sOutput:\n%s\n", i+1, tc.desc, err, tc.input, out)
			os.Exit(1)
		}
		for idx := 0; idx < tc.t; idx++ {
			if exp[idx] != got[idx] {
				fmt.Printf("Wrong answer on test %d (%s) case %d: expected %d got %d\nInput:\n%s", i+1, tc.desc, idx+1, exp[idx], got[idx], tc.input)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	path := "./ref2161F.bin"
	cmd := exec.Command("go", "build", "-o", path, "2161F.go")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	return path, nil
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
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func parseAnswers(out string, t int) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	ans := make([]int64, t)
	for i := 0; i < t; i++ {
		if _, err := fmt.Fscan(reader, &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to read answer %d: %v", i+1, err)
		}
	}
	if extra := strings.TrimSpace(readRemaining(reader)); extra != "" {
		return nil, fmt.Errorf("unexpected extra output: %q", extra)
	}
	return ans, nil
}

func readRemaining(r *bufio.Reader) string {
	var sb strings.Builder
	for {
		line, err := r.ReadString('\n')
		sb.WriteString(line)
		if err != nil {
			break
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	add := func(desc string, cases []treeCase) {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", len(cases))
		for _, c := range cases {
			fmt.Fprintf(&sb, "%d\n", c.n)
			for _, e := range c.edges {
				fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
			}
		}
		tests = append(tests, testCase{
			desc:  desc,
			input: sb.String(),
			t:     len(cases),
		})
	}

	add("trivial", []treeCase{
		{n: 1},
		{n: 2, edges: [][2]int{{1, 2}}},
	})

	add("line-small", []treeCase{
		buildLine(5),
		buildStar(6),
	})

	rng := rand.New(rand.NewSource(123456789))
	for len(tests) < 60 {
		numCases := rng.Intn(4) + 1
		cases := make([]treeCase, numCases)
		for i := 0; i < numCases; i++ {
			n := rng.Intn(40) + 1
			cases[i] = buildRandomTree(n, rng)
		}
		add(fmt.Sprintf("random-small-%d", len(tests)), cases)
	}

	// Stress tests with larger n
	large1 := buildLine(5000)
	large2 := buildRandomTree(5000, rng)
	large3 := buildStar(5000)
	add("large", []treeCase{large1, large2, large3})

	return tests
}

type treeCase struct {
	n     int
	edges [][2]int
}

func buildLine(n int) treeCase {
	edges := make([][2]int, 0, n-1)
	for i := 1; i < n; i++ {
		edges = append(edges, [2]int{i, i + 1})
	}
	return treeCase{n: n, edges: edges}
}

func buildStar(n int) treeCase {
	if n == 0 {
		return treeCase{n: 0}
	}
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return treeCase{n: n, edges: edges}
}

func buildRandomTree(n int, rng *rand.Rand) treeCase {
	if n <= 1 {
		return treeCase{n: n}
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	return treeCase{n: n, edges: edges}
}
