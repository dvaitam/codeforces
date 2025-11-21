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

type testCase struct {
	desc  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, tc := range tests {
		expOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		gotOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d (%s): %v\ninput:\n%s", i+1, tc.desc, err, tc.input)
			os.Exit(1)
		}
		expTokens := normalizeTokens(expOut)
		gotTokens := normalizeTokens(gotOut)
		if len(expTokens) != len(gotTokens) {
			fmt.Fprintf(os.Stderr, "test %d (%s): token count mismatch\nexpected: %v\ngot: %v\n", i+1, tc.desc, expTokens, gotTokens)
			os.Exit(1)
		}
		for idx := range expTokens {
			if expTokens[idx] != gotTokens[idx] {
				fmt.Fprintf(os.Stderr, "test %d (%s): mismatch at answer %d\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, tc.desc, idx+1, tc.input, expTokens[idx], gotTokens[idx])
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1592C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	src := filepath.Clean("1000-1999/1500-1599/1590-1599/1592/1592C.go")
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runProgram(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func normalizeTokens(out string) []string {
	fields := strings.Fields(out)
	res := make([]string, len(fields))
	for i, f := range fields {
		res[i] = strings.ToUpper(f)
	}
	return res
}

func buildTests() []testCase {
	tests := []testCase{
		{desc: "total_zero", input: formatSingleCase(1, simpleCase(2, 2, []int{1, 1}, [][]int{{1, 2}}))},
		{desc: "k_too_small", input: formatSingleCase(1, simpleCase(2, 2, []int{1, 2}, [][]int{{1, 2}}))},
		{desc: "line_three", input: formatSingleCase(1, simpleCase(3, 3, []int{1, 1, 0}, [][]int{{1, 2}, {2, 3}}))},
		{desc: "star_five", input: formatSingleCase(1, simpleCase(5, 4, []int{1, 2, 3, 4, 5}, [][]int{{1, 2}, {1, 3}, {1, 4}, {4, 5}}))},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomMultiCase(fmt.Sprintf("rand-%d", i+1), rng))
	}
	return tests
}

type caseData struct {
	n     int
	k     int
	a     []int
	edges [][]int
}

func simpleCase(n, k int, a []int, edges [][]int) caseData {
	return caseData{n: n, k: k, a: a, edges: edges}
}

func formatSingleCase(t int, c caseData) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	appendCase(&sb, c)
	return sb.String()
}

func appendCase(sb *strings.Builder, c caseData) {
	sb.WriteString(fmt.Sprintf("%d %d\n", c.n, c.k))
	for i, v := range c.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for _, e := range c.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
}

func randomMultiCase(desc string, rng *rand.Rand) testCase {
	t := rng.Intn(4) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	totalN := 0
	for i := 0; i < t; i++ {
		n := rng.Intn(15) + 2
		k := rng.Intn(n-1) + 2
		totalN += n
		if totalN > 60 {
			n = 2 + rng.Intn(6)
			k = rng.Intn(n-1) + 2
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(1_000_000_000) + 1
		}
		edges := randomTree(n, rng)
		appendCase(&sb, caseData{n: n, k: k, a: a, edges: edges})
	}
	return testCase{desc: desc, input: sb.String()}
}

func randomTree(n int, rng *rand.Rand) [][]int {
	edges := make([][]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, []int{u, v})
	}
	return edges
}
