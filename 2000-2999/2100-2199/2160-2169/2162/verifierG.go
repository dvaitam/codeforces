package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// refSourceG points to the local reference solution to avoid GOPATH resolution.
const refSourceG = "2162G.go"

type testCaseG struct {
	n int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReferenceG()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTestsG()
	input := buildInputG(tests)

	refOut, err := runProgramG(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidateG(os.Args[1], input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseAnswersG(refOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseAnswersG(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		refSol := refAns[idx]
		candSol := candAns[idx]

		if refSol.exists != candSol.exists {
			fmt.Fprintf(os.Stderr, "test %d: existence mismatch (reference %v, candidate %v)\n", idx+1, refSol.exists, candSol.exists)
			os.Exit(1)
		}

		if !candSol.exists {
			continue
		}

		if err := validateTree(tc.n, candSol.edges); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: invalid tree: %v\n", idx+1, err)
			os.Exit(1)
		}

		candSum := sumProducts(candSol.edges)
		refSum := sumProducts(refSol.edges)
		if !isPerfectSquare(candSum) {
			fmt.Fprintf(os.Stderr, "test %d: candidate sum %d is not a perfect square\n", idx+1, candSum)
			os.Exit(1)
		}
		if candSum != refSum {
			fmt.Fprintf(os.Stderr, "test %d: candidate sum %d differs from reference sum %d\n", idx+1, candSum, refSum)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceG() (string, error) {
	tmp, err := os.CreateTemp("", "2162G-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceG))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runCandidateG(path, input string) (string, error) {
	cmd := commandForG(path)
	return runWithInputG(cmd, input)
}

func runProgramG(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInputG(cmd, input)
}

func commandForG(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInputG(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

type answerG struct {
	exists bool
	edges  [][2]int
}

func parseAnswersG(output string, tests []testCaseG) ([]answerG, error) {
	tokens := strings.Fields(output)
	ans := make([]answerG, len(tests))
	pos := 0
	for i, tc := range tests {
		if pos >= len(tokens) {
			return nil, fmt.Errorf("missing output for test %d", i+1)
		}
		if tokens[pos] == "-1" {
			ans[i] = answerG{}
			pos++
			continue
		}
		need := tc.n - 1
		edges := make([][2]int, 0, need)
		for j := 0; j < need; j++ {
			if pos+1 >= len(tokens) {
				return nil, fmt.Errorf("test %d: insufficient tokens for edge %d", i+1, j+1)
			}
			u, err := strconv.Atoi(tokens[pos])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, tokens[pos])
			}
			v, err := strconv.Atoi(tokens[pos+1])
			if err != nil {
				return nil, fmt.Errorf("test %d: invalid integer %q", i+1, tokens[pos+1])
			}
			edges = append(edges, [2]int{u, v})
			pos += 2
		}
		ans[i] = answerG{exists: true, edges: edges}
	}
	if pos != len(tokens) {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[pos])
	}
	return ans, nil
}

func validateTree(n int, edges [][2]int) error {
	if len(edges) != n-1 {
		return fmt.Errorf("expected %d edges, got %d", n-1, len(edges))
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("edge (%d,%d) has invalid vertex", u, v)
		}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	visited := make([]bool, n+1)
	var stack []int
	stack = append(stack, 1)
	visited[1] = true
	count := 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		for _, to := range adj[v] {
			if !visited[to] {
				visited[to] = true
				stack = append(stack, to)
			}
		}
	}
	if count != n {
		return fmt.Errorf("graph is disconnected")
	}
	return nil
}

func sumProducts(edges [][2]int) int64 {
	var total int64
	for _, e := range edges {
		total += int64(e[0]) * int64(e[1])
	}
	return total
}

func isPerfectSquare(x int64) bool {
	if x < 0 {
		return false
	}
	root := int64(math.Sqrt(float64(x)))
	for (root+1)*(root+1) <= x {
		root++
	}
	for root*root > x {
		root--
	}
	return root*root == x
}

func generateTestsG() []testCaseG {
	const limit = 200000
	rng := rand.New(rand.NewSource(21622062))
	var tests []testCaseG
	total := 0

	add := func(n int) {
		if n < 2 {
			n = 2
		}
		if total+n > limit {
			return
		}
		tests = append(tests, testCaseG{n: n})
		total += n
	}

	for _, n := range []int{2, 3, 4, 5, 6, 7, 8, 9, 10} {
		add(n)
	}

	for i := 0; i < 60 && total+2 <= limit; i++ {
		add(2 + rng.Intn(25))
	}
	for i := 0; i < 40 && total+30 <= limit; i++ {
		add(30 + rng.Intn(200))
	}
	for i := 0; i < 20 && total+200 <= limit; i++ {
		add(200 + rng.Intn(2000))
	}
	for i := 0; i < 10 && total+2000 <= limit; i++ {
		add(2000 + rng.Intn(5000))
	}

	if total+50000 <= limit {
		add(50000)
	}
	if limit-total >= 2 {
		add(limit - total)
	}

	return tests
}

func buildInputG(tests []testCaseG) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
	}
	return b.String()
}
