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
	name  string
	input string
}

func solveRef(input string) ([]int64, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return nil, err
	}
	colors := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &colors[i])
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	st := make([]int, n+1)
	en := make([]int, n+1)
	flat := make([]int, n+1)
	sz := make([]int, n+1)
	heavy := make([]int, n+1)
	timer := 0
	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		timer++
		st[u] = timer
		flat[timer] = u
		sz[u] = 1
		heavy[u] = -1
		maxSize := 0
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs1(v, u)
			if sz[v] > maxSize {
				maxSize = sz[v]
				heavy[u] = v
			}
			sz[u] += sz[v]
		}
		en[u] = timer
	}
	dfs1(1, 0)

	freqColor := make([]int, n+1)
	cntFreq := make([]int, n+1)
	sumFreq := make([]int64, n+1)
	ans := make([]int64, n+1)
	maxFreq := 0

	addNode := func(u int) {
		c := colors[u]
		old := freqColor[c]
		if old > 0 {
			cntFreq[old]--
			sumFreq[old] -= int64(c)
		}
		newf := old + 1
		freqColor[c] = newf
		cntFreq[newf]++
		sumFreq[newf] += int64(c)
		if newf > maxFreq {
			maxFreq = newf
		}
	}

	removeNode := func(u int) {
		c := colors[u]
		f := freqColor[c]
		if f == 0 {
			return
		}
		cntFreq[f]--
		sumFreq[f] -= int64(c)
		freqColor[c] = 0
		if cntFreq[maxFreq] == 0 {
			for maxFreq > 0 && cntFreq[maxFreq] == 0 {
				maxFreq--
			}
		}
	}

	var dfs2 func(u, p int, keep bool)
	dfs2 = func(u, p int, keep bool) {
		for _, v := range adj[u] {
			if v == p || v == heavy[u] {
				continue
			}
			dfs2(v, u, false)
		}
		if heavy[u] != -1 {
			dfs2(heavy[u], u, true)
		}
		for _, v := range adj[u] {
			if v == p || v == heavy[u] {
				continue
			}
			for i := st[v]; i <= en[v]; i++ {
				addNode(flat[i])
			}
		}
		addNode(u)
		ans[u] = sumFreq[maxFreq]
		if !keep {
			for i := st[u]; i <= en[u]; i++ {
				removeNode(flat[i])
			}
			maxFreq = 0
		}
	}
	dfs2(1, 0, true)
	return ans[1:], nil
}

func makeCase(name string, n int, colors []int, edges [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", colors[i+1]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return testCase{name: name, input: sb.String()}
}

func randomTree(rng *rand.Rand, n int) [][2]int {
	edges := make([][2]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges[i-2] = [2]int{p, i}
	}
	return edges
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(600))
	var tests []testCase
	gen := func(prefix string, count, maxN int) {
		for i := 0; i < count; i++ {
			n := rng.Intn(maxN-1) + 1
			if n == 0 {
				n = 1
			}
			colors := make([]int, n+1)
			for j := 1; j <= n; j++ {
				colors[j] = rng.Intn(n) + 1
			}
			edges := randomTree(rng, n)
			tests = append(tests, makeCase(fmt.Sprintf("%s_%d", prefix, i+1), n, colors, edges))
		}
	}
	gen("small", 80, 10)
	gen("medium", 80, 50)
	gen("large", 40, 200)
	return tests
}

func handcraftedTests() []testCase {
	return []testCase{
		makeCase("single_node", 1, []int{0, 5}, [][2]int{}),
		makeCase("line_same_color", 4, []int{0, 1, 1, 1, 1}, [][2]int{{1, 2}, {2, 3}, {3, 4}}),
		makeCase("line_diff_colors", 4, []int{0, 1, 2, 3, 4}, [][2]int{{1, 2}, {2, 3}, {3, 4}}),
		makeCase("star", 5, []int{0, 1, 2, 3, 4, 5}, [][2]int{{1, 2}, {1, 3}, {1, 4}, {1, 5}}),
		makeCase("balanced", 7, []int{0, 1, 2, 2, 3, 3, 3, 3}, [][2]int{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {3, 7}}),
	}
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseOutput(output string, n int) ([]int64, error) {
	fields := strings.Fields(output)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d integers, got %d", n, len(fields))
	}
	ans := make([]int64, n)
	for i, f := range fields {
		if _, err := fmt.Sscan(f, &ans[i]); err != nil {
			return nil, fmt.Errorf("failed to parse integer #%d: %v", i+1, err)
		}
	}
	return ans, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := append(handcraftedTests(), randomTests()...)
	for idx, tc := range tests {
		expect, err := solveRef(tc.input)
		if err != nil {
			fmt.Printf("failed to compute reference for %s: %v\n", tc.name, err)
			os.Exit(1)
		}
		output, runErr := runCandidate(bin, tc.input)
		if runErr != nil {
			fmt.Printf("test %d (%s) runtime error: %v\ninput:\n%s", idx+1, tc.name, runErr, tc.input)
			os.Exit(1)
		}
		got, parseErr := parseOutput(output, len(expect))
		if parseErr != nil {
			fmt.Printf("test %d (%s) invalid output: %v\ninput:\n%s\noutput:\n%s\n", idx+1, tc.name, parseErr, tc.input, output)
			os.Exit(1)
		}
		match := true
		for i := range expect {
			if expect[i] != got[i] {
				match = false
				break
			}
		}
		if !match {
			fmt.Printf("test %d (%s) mismatch\ninput:\n%s\nexpect:%v\nactual:%v\n", idx+1, tc.name, tc.input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
