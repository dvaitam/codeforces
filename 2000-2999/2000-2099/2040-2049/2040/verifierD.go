package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2040D.go"

type testCase struct {
	name  string
	n     int
	edges [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate_binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	maxN := 0
	for _, tc := range tests {
		if tc.n > maxN {
			maxN = tc.n
		}
	}
	primes := sieve(2 * maxN)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference solution: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	refHas, err := checkOutputs(tests, refOut, primes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference produced invalid output: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}
	candHas, err := checkOutputs(tests, candOut, primes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}

	for i, hasSol := range refHas {
		if hasSol && !candHas[i] {
			fmt.Fprintf(os.Stderr, "test %d (%s): candidate printed -1 but a valid solution exists\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				i+1, tests[i].name, input, refOut, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier directory")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2040D-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("reference build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
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

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 32)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
	}
	return sb.String()
}

func checkOutputs(tests []testCase, output string, primes []bool) ([]bool, error) {
	tokens := strings.Fields(output)
	idx := 0
	hasSol := make([]bool, len(tests))
	for caseIdx, tc := range tests {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("case %d: missing output", caseIdx+1)
		}
		tok := tokens[idx]
		idx++
		if tok == "-1" {
			hasSol[caseIdx] = false
			continue
		}
		arr := make([]int, tc.n)
		used := make([]bool, 2*tc.n+1)
		val, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("case %d: invalid integer %q", caseIdx+1, tok)
		}
		if val < 1 || val > 2*tc.n {
			return nil, fmt.Errorf("case %d: value %d out of range 1..%d", caseIdx+1, val, 2*tc.n)
		}
		used[val] = true
		arr[0] = val
		for i := 1; i < tc.n; i++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("case %d: expected %d numbers, got %d", caseIdx+1, tc.n, i)
			}
			tok = tokens[idx]
			idx++
			val, err = strconv.Atoi(tok)
			if err != nil {
				return nil, fmt.Errorf("case %d: invalid integer %q", caseIdx+1, tok)
			}
			if val < 1 || val > 2*tc.n {
				return nil, fmt.Errorf("case %d: value %d out of range 1..%d", caseIdx+1, val, 2*tc.n)
			}
			if used[val] {
				return nil, fmt.Errorf("case %d: duplicate value %d", caseIdx+1, val)
			}
			used[val] = true
			arr[i] = val
		}
		for _, e := range tc.edges {
			u := e[0] - 1
			v := e[1] - 1
			diff := abs(arr[u] - arr[v])
			if diff == 0 {
				return nil, fmt.Errorf("case %d: zero difference on edge %d-%d", caseIdx+1, e[0], e[1])
			}
			if diff < len(primes) && primes[diff] {
				return nil, fmt.Errorf("case %d: prime difference %d on edge %d-%d", caseIdx+1, diff, e[0], e[1])
			}
		}
		hasSol[caseIdx] = true
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("expected %d tokens, got %d", idx, len(tokens))
	}
	return hasSol, nil
}

func buildTests() []testCase {
	tests := []testCase{
		newPathCase("two_nodes", 2),
		newPathCase("short_path", 4),
		newStarCase("small_star", 6),
		newPathCase("prime_len_path", 11),
		newRandomCase("balanced_tree", 15, rand.New(rand.NewSource(7))),
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	// A few medium structured cases
	tests = append(tests, newPathCase("path_2000", 2000))
	tests = append(tests, newStarCase("star_1500", 1500))
	totalN += 3500

	// Random small cases
	for i := 0; i < 50 && totalN < 50000; i++ {
		n := rng.Intn(40) + 2
		tests = append(tests, newRandomCase(fmt.Sprintf("small_rand_%d", i+1), n, rng))
		totalN += n
	}

	// Random medium cases
	for i := 0; i < 25 && totalN < 120000; i++ {
		n := rng.Intn(4000) + 500
		if totalN+n > 200000 {
			break
		}
		tests = append(tests, newRandomCase(fmt.Sprintf("mid_rand_%d", i+1), n, rng))
		totalN += n
	}

	// One or two larger cases close to the limit
	if totalN < 200000 {
		n := 80000
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		if n >= 2 {
			tests = append(tests, newRandomCase("large_rand", n, rng))
			totalN += n
		}
	}
	if totalN < 200000 {
		n := 200000 - totalN
		if n >= 2 {
			tests = append(tests, newPathCase("large_path", n))
		}
	}

	return tests
}

func newPathCase(name string, n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{i - 1, i})
	}
	return testCase{name: name, n: n, edges: edges}
}

func newStarCase(name string, n int) testCase {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		edges = append(edges, [2]int{1, i})
	}
	return testCase{name: name, n: n, edges: edges}
}

func newRandomCase(name string, n int, rng *rand.Rand) testCase {
	edges := make([][2]int, 0, n-1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		edges = append(edges, [2]int{p, v})
	}
	return testCase{name: name, n: n, edges: edges}
}

func sieve(limit int) []bool {
	isPrime := make([]bool, limit+1)
	if limit < 2 {
		return isPrime
	}
	for i := 2; i <= limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i <= limit; i++ {
		if isPrime[i] {
			for j := i * i; j <= limit; j += i {
				isPrime[j] = false
			}
		}
	}
	return isPrime
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
