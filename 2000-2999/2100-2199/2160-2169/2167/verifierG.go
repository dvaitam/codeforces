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
)

const refSourceG = "2000-2999/2100-2199/2160-2169/2167/2167G.go"

type testCaseG struct {
	n int
	a []int
	c []int
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

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference execution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	candOut, err := runCandidate(os.Args[1], input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate execution failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		expected := refAns[i]
		got := candAns[i]
		if expected != got {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\n", i+1, expected, got)
			os.Exit(1)
		}
		if !checkFeasibility(tc, got) {
			fmt.Fprintf(os.Stderr, "test %d: candidate answer %d is infeasible\n", i+1, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReferenceG() (string, error) {
	tmp, err := os.CreateTemp("", "2167G-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(output string, expected int) ([]int64, error) {
	tokens := strings.Fields(output)
	if len(tokens) < expected {
		return nil, fmt.Errorf("expected %d numbers, got %d", expected, len(tokens))
	}
	if len(tokens) > expected {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[expected])
	}
	ans := make([]int64, expected)
	for i := 0; i < expected; i++ {
		val, err := strconv.ParseInt(tokens[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", tokens[i])
		}
		ans[i] = val
	}
	return ans, nil
}

func checkFeasibility(tc testCaseG, target int64) bool {
	best := solveDP(tc.a, tc.c)
	return best == target
}

func solveDP(a, c []int) int64 {
	n := len(a)
	if n == 0 {
		return 0
	}
	coords := make([]pair, n)
	for i, v := range a {
		coords[i] = pair{val: v, pos: i}
	}
	sortPairs(coords)
	comp := make([]int, n)
	last := -1
	ptr := 0
	for _, p := range coords {
		if last != p.val {
			ptr++
			last = p.val
		}
		comp[p.pos] = ptr
	}
	bit := newBIT(ptr + 2)
	var total int64
	for _, v := range c {
		total += int64(v)
	}
	var best int64
	for i := 0; i < n; i++ {
		idx := comp[i]
		cur := bit.query(idx) + int64(c[i])
		if cur > best {
			best = cur
		}
		bit.update(idx, cur)
	}
	return total - best
}

type bitTree struct {
	tree []int64
}

func newBIT(size int) *bitTree {
	return &bitTree{tree: make([]int64, size)}
}

func (b *bitTree) update(idx int, val int64) {
	for idx < len(b.tree) {
		if val > b.tree[idx] {
			b.tree[idx] = val
		}
		idx += idx & -idx
	}
}

func (b *bitTree) query(idx int) int64 {
	var res int64
	for idx > 0 {
		if b.tree[idx] > res {
			res = b.tree[idx]
		}
		idx -= idx & -idx
	}
	return res
}

func generateTestsG() []testCaseG {
	rng := rand.New(rand.NewSource(21672077))
	var tests []testCaseG
	total := 0
	const limit = 8000

	add := func(a, c []int) {
		if len(a) != len(c) || len(a) == 0 {
			return
		}
		if total+len(a) > limit {
			return
		}
		tests = append(tests, testCaseG{
			n: len(a),
			a: append([]int(nil), a...),
			c: append([]int(nil), c...),
		})
		total += len(a)
	}

	add([]int{1}, []int{5})
	add([]int{1, 1}, []int{3, 4})
	add([]int{1, 2, 2}, []int{2, 2, 2})
	add([]int{3, 2, 1}, []int{1, 10, 100})
	add([]int{5, 4, 3, 2, 1}, []int{1, 1, 1, 1, 1})

	for i := 0; i < 50; i++ {
		n := 1 + rng.Intn(10)
		add(randomArray(rng, n), randomCosts(rng, n))
	}

	for i := 0; i < 40; i++ {
		n := 10 + rng.Intn(50)
		add(randomArray(rng, n), randomCosts(rng, n))
	}

	for total < limit {
		n := min(limit-total, 200)
		add(randomArray(rng, n), randomCosts(rng, n))
	}

	return tests
}

func randomArray(rng *rand.Rand, n int) []int {
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(1_000_000_000) + 1
	}
	return arr
}

func randomCosts(rng *rand.Rand, n int) []int {
	c := make([]int, n)
	for i := range c {
		c[i] = rng.Intn(1_000_000_000) + 1
	}
	return c
}

func buildInputG(tests []testCaseG) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d\n", tc.n)
		writeInts(&b, tc.a)
		writeInts(&b, tc.c)
	}
	return b.String()
}

func writeInts(b *strings.Builder, arr []int) {
	for i, v := range arr {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}

type pair struct {
	val int
	pos int
}

func sortPairs(arr []pair) {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j].val > key.val {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
