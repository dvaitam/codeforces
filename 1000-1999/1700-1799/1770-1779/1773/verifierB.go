package main

import (
	"bufio"
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

type testCase struct {
	input string
	n     int
	k     int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for idx, tc := range tests {
		refOut, err := runProgram(ref, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		refParents, err := parseParents(refOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}
		if err := validateTree(refParents, tc.n); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid tree on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotParents, err := parseParents(gotOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}
		if err := validateTree(gotParents, tc.n); err != nil {
			fmt.Fprintf(os.Stderr, "participant produced invalid tree on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if ok, err := compareTrees(tc.input, gotParents, tc.n, tc.k); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to verify records: %v\ninput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		} else if !ok {
			fmt.Fprintf(os.Stderr, "test %d: participant tree inconsistent with journal\ninput:\n%sreference output:\n%s\nparticipant output:\n%s\n", idx+1, tc.input, refOut, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	dir, err := verifierDir()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "1773B_ref_*.bin")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "1773B.go")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func verifierDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("unable to determine verifier directory")
	}
	return filepath.Dir(file), nil
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseParents(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers got %d", n, len(fields))
	}
	parent := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val, err := strconv.Atoi(fields[i-1])
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", fields[i-1])
		}
		parent[i] = val
	}
	return parent, nil
}

func validateTree(parent []int, n int) error {
	root := -1
	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if parent[i] == -1 {
			if root != -1 {
				return fmt.Errorf("multiple roots found")
			}
			root = i
			continue
		}
		if parent[i] < 1 || parent[i] > n {
			return fmt.Errorf("invalid parent for node %d: %d", i, parent[i])
		}
		children[parent[i]] = append(children[parent[i]], i)
	}
	if root == -1 {
		return fmt.Errorf("no root found")
	}
	for i := 1; i <= n; i++ {
		if len(children[i]) != 0 && len(children[i]) != 2 {
			return fmt.Errorf("node %d has %d children (must be 0 or 2)", i, len(children[i]))
		}
	}
	visited := make([]bool, n+1)
	stack := []int{root}
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if visited[v] {
			continue
		}
		visited[v] = true
		stack = append(stack, children[v]...)
	}
	for i := 1; i <= n; i++ {
		if !visited[i] {
			return fmt.Errorf("node %d unreachable from root", i)
		}
	}
	return nil
}

func compareTrees(input string, parent []int, n, k int) (bool, error) {
	parsedN, parsedK, records, err := parseRecords(input)
	if err != nil {
		return false, err
	}
	if parsedN != n || parsedK != k {
		return false, fmt.Errorf("input mismatch with test metadata")
	}
	for _, seq := range records {
		if !checkSequence(n, parent, seq) {
			return false, nil
		}
	}
	return true, nil
}

func parseRecords(input string) (int, int, [][]int, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return 0, 0, nil, fmt.Errorf("failed to read n,k: %v", err)
	}
	records := make([][]int, k)
	for i := 0; i < k; i++ {
		records[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(reader, &records[i][j]); err != nil {
				return 0, 0, nil, fmt.Errorf("failed to read record: %v", err)
			}
		}
	}
	return n, k, records, nil
}

func checkSequence(n int, parent []int, seq []int) bool {
	children := make([][]int, n+1)
	root := -1
	for i := 1; i <= n; i++ {
		if parent[i] == -1 {
			root = i
			continue
		}
		children[parent[i]] = append(children[parent[i]], i)
	}
	if root == -1 {
		return false
	}
	pos := make([]int, n+1)
	for idx, v := range seq {
		if v < 1 || v > n {
			return false
		}
		pos[v] = idx
	}
	stack := []int{root}
	idx := 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if idx >= n || seq[idx] != v {
			return false
		}
		idx++
		if len(children[v]) == 0 {
			continue
		}
		if len(children[v]) != 2 {
			return false
		}
		c1, c2 := children[v][0], children[v][1]
		p1, p2 := pos[c1], pos[c2]
		if p1 <= idx-1 || p2 <= idx-1 {
			return false
		}
		if p1 < p2 {
			stack = append(stack, c2, c1)
		} else {
			stack = append(stack, c1, c2)
		}
	}
	return idx == n
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 40)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	rng := rand.New(rand.NewSource(123))
	return []testCase{
		makeTestCase(3, 60, rng),
		makeTestCase(5, 70, rng),
	}
}

func randomTests(rng *rand.Rand, batches int) []testCase {
	var tests []testCase
	for i := 0; i < batches; i++ {
		n := rng.Intn(31) + 3
		if n%2 == 0 {
			n++
		}
		if n > 999 {
			n = 999
		}
		k := rng.Intn(51) + 50
		tests = append(tests, makeTestCase(n, k, rng))
	}
	return tests
}

func stressTests() []testCase {
	rng := rand.New(rand.NewSource(321))
	return []testCase{
		makeTestCase(999, 100, rng),
	}
}

func makeTestCase(n, k int, rng *rand.Rand) testCase {
	if n%2 == 0 {
		n++
	}
	if n > 999 {
		n = 999
	}
	nodes := make([]int, n)
	for i := 0; i < n; i++ {
		nodes[i] = i + 1
	}
	rng.Shuffle(n, func(i, j int) {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	})
	parent := make([]int, n+1)
	root := buildTree(nodes, 0, parent, rng)
	parent[root] = -1

	children := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		if parent[i] > 0 {
			children[parent[i]] = append(children[parent[i]], i)
		}
	}
	records := make([][]int, k)
	for i := 0; i < k; i++ {
		records[i] = randomRecord(root, children, rng)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for _, rec := range records {
		for j, val := range rec {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String(), n: n, k: k}
}

func buildTree(nodes []int, parent int, parents []int, rng *rand.Rand) int {
	if len(nodes) == 1 {
		parents[nodes[0]] = parent
		return nodes[0]
	}
	root := nodes[0]
	parents[root] = parent
	remaining := nodes[1:]
	left := randomOdd(rng, len(remaining)-1)
	leftNodes := append([]int(nil), remaining[:left]...)
	rightNodes := append([]int(nil), remaining[left:]...)
	buildTree(leftNodes, root, parents, rng)
	buildTree(rightNodes, root, parents, rng)
	return root
}

func randomOdd(rng *rand.Rand, limit int) int {
	if limit <= 0 {
		return 1
	}
	for {
		val := rng.Intn(limit) + 1
		if val%2 == 1 {
			return val
		}
	}
}

func randomRecord(root int, children [][]int, rng *rand.Rand) []int {
	res := make([]int, 0)
	var dfs func(int)
	dfs = func(v int) {
		res = append(res, v)
		if len(children[v]) == 0 {
			return
		}
		if rng.Intn(2) == 0 {
			dfs(children[v][0])
			dfs(children[v][1])
		} else {
			dfs(children[v][1])
			dfs(children[v][0])
		}
	}
	dfs(root)
	return res
}
