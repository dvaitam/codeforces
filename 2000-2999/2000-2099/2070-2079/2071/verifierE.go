package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const mod = 998244353

type frac struct {
	p int64
	q int64
}

type testCase struct {
	n     int
	fracs []frac
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := buildTestCases(rng)
	input := buildInput(cases)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	expected, err := parseAnswers(refOut, len(cases))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid reference output: %v\n", err)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\n", err)
		os.Exit(1)
	}
	got, err := parseAnswers(candOut, len(cases))
	if err != nil {
		fmt.Fprintf(os.Stderr, "bad candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\ninput:\n%s", i+1, expected[i], got[i], input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	const refName = "./ref_2071E.bin"
	cmd := exec.Command("go", "build", "-o", refName, "2071E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return refName, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return stdout.String(), nil
}

func parseAnswers(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d answers, got %d (output: %q)", expected, len(fields), out)
	}
	res := make([]int64, expected)
	for i, f := range fields {
		val, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q: %v", f, err)
		}
		if val < 0 || val >= mod {
			return nil, fmt.Errorf("answer %d out of range [0, %d)", val, mod)
		}
		res[i] = val
	}
	return res, nil
}

func buildTestCases(rng *rand.Rand) []testCase {
	cases := deterministicCases()
	totalN := 0
	for _, tc := range cases {
		totalN += tc.n
	}

	for len(cases) < 150 && totalN < 100000 {
		n := rng.Intn(4000) + 1
		if totalN+n > 100000 {
			n = 100000 - totalN
		}
		if n <= 0 {
			break
		}
		cases = append(cases, randomCase(rng, n))
		totalN += n
	}
	return cases
}

func deterministicCases() []testCase {
	return []testCase{
		{
			n:     1,
			fracs: []frac{{p: 1, q: 2}},
			edges: nil,
		},
		{
			n:     2,
			fracs: []frac{{1, 3}, {2, 3}},
			edges: [][2]int{{1, 2}},
		},
		{
			n: 3,
			fracs: []frac{
				{1, 4},
				{2, 5},
				{3, 7},
			},
			edges: [][2]int{{1, 2}, {2, 3}},
		},
		{
			n: 5,
			fracs: []frac{
				{1, 9}, {2, 9}, {3, 9}, {4, 9}, {5, 9},
			},
			edges: [][2]int{{1, 2}, {1, 3}, {1, 4}, {4, 5}},
		},
	}
}

func randomCase(rng *rand.Rand, n int) testCase {
	fracs := make([]frac, n)
	for i := 0; i < n; i++ {
		fracs[i] = randomFrac(rng)
	}
	edges := randomTree(rng, n)
	return testCase{n: n, fracs: fracs, edges: edges}
}

func randomFrac(rng *rand.Rand) frac {
	q := int64(rng.Intn(mod-2) + 2) // [2, mod-1]
	p := int64(rng.Intn(int(q-1)) + 1)
	return frac{p: p, q: q}
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	if n <= 1 {
		return nil
	}
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{u, v})
	}
	rng.Shuffle(len(edges), func(i, j int) {
		edges[i], edges[j] = edges[j], edges[i]
	})
	return edges
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, tc := range cases {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, f := range tc.fracs {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(fmt.Sprintf("%d %d", f.p, f.q))
		}
		sb.WriteByte('\n')
		for i, e := range tc.edges {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(fmt.Sprintf("%d %d", e[0], e[1]))
		}
		if len(tc.edges) > 0 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}
