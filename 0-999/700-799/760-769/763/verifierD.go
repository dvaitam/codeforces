package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type edge struct {
	u, v int
}

type testCase struct {
	input []byte
	n     int
	adj   [][]int
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	out := filepath.Join(dir, "oracleD")
	cmd := exec.Command("go", "build", "-o", out, "763D.go")
	if outBytes, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build oracle: %v\n%s", err, string(outBytes))
	}
	return out, nil
}

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseAnswer(out string, n int) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 1 || val > n {
		return 0, fmt.Errorf("vertex %d out of range 1..%d", val, n)
	}
	return val, nil
}

func makeTest(n int, edges []edge) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	adj := make([][]int, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}
	return testCase{
		input: []byte(sb.String()),
		n:     n,
		adj:   adj,
	}
}

func deterministicTests() []testCase {
	tests := []testCase{
		makeTest(1, nil),
		makeTest(2, []edge{{1, 2}}),
		makeTest(3, []edge{{1, 2}, {2, 3}}),
		makeTest(5, []edge{{1, 2}, {1, 3}, {1, 4}, {1, 5}}),
		makeTest(7, []edge{{1, 2}, {1, 3}, {2, 4}, {2, 5}, {3, 6}, {3, 7}}),
		makeTest(9, []edge{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}}),
	}
	return tests
}

func randomTreeEdges(n int, rnd *rand.Rand) []edge {
	if n <= 1 {
		return nil
	}
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := rnd.Intn(i-1) + 1
		edges = append(edges, edge{u: i, v: parent})
	}
	return edges
}

func randomTests(count int) []testCase {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, count+1)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%12 == 0:
			n = rnd.Intn(90000) + 10000
		case i%4 == 0:
			n = rnd.Intn(4000) + 200
		default:
			n = rnd.Intn(200) + 2
		}
		edges := randomTreeEdges(n, rnd)
		tests = append(tests, makeTest(n, edges))
	}
	tests = append(tests, makeTest(100000, randomTreeEdges(100000, rand.New(rand.NewSource(time.Now().UnixNano()+1)))))
	return tests
}

func encode(vals []int) string {
	if len(vals) == 0 {
		return "#"
	}
	var sb strings.Builder
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func evaluateTree(adj [][]int) (int, []int) {
	n := len(adj) - 1
	if n <= 0 {
		return 0, nil
	}
	sub := make([]int, n+1)
	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		sub[u] = 1
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs1(v, u)
			sub[u] += sub[v]
		}
	}
	dfs1(1, 0)
	centroid := 1
	bestSub := n + 1
	for i := 1; i <= n; i++ {
		if sub[i]*2 > n && sub[i] < bestSub {
			bestSub = sub[i]
			centroid = i
		}
	}
	children := make([][]int, n+1)
	height := make([]int, n+1)
	var dfs2 func(u, p int)
	dfs2 = func(u, p int) {
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			children[u] = append(children[u], v)
			dfs2(v, u)
			if height[v]+1 > height[u] {
				height[u] = height[v] + 1
			}
		}
	}
	dfs2(centroid, 0)
	maxDepth := 0
	for i := 1; i <= n; i++ {
		if height[i] > maxDepth {
			maxDepth = height[i]
		}
	}
	layers := make([][]int, maxDepth+1)
	for i := 1; i <= n; i++ {
		layers[height[i]] = append(layers[height[i]], i)
	}
	num := make([]int, n+1)
	cnt := 0
	for d := 0; d <= maxDepth; d++ {
		if len(layers[d]) == 0 {
			continue
		}
		shapeMap := make(map[string]int)
		for _, u := range layers[d] {
			tp := make([]int, len(children[u]))
			for idx, v := range children[u] {
				tp[idx] = num[v]
			}
			if len(tp) > 1 {
				sort.Ints(tp)
			}
			key := encode(tp)
			if id, ok := shapeMap[key]; ok {
				num[u] = id
			} else {
				cnt++
				shapeMap[key] = cnt
				num[u] = cnt
			}
		}
	}
	if cnt == 0 {
		cnt = 1
	}
	S := make([]int, cnt+1)
	for i := 1; i <= n; i++ {
		S[num[i]]++
	}
	R := 0
	for i := 1; i <= cnt; i++ {
		if S[i] > 0 {
			R++
		}
	}
	scores := make([]int, n+1)
	best := -1 << 60
	var dfsScore func(u, depth int)
	dfsScore = func(u, depth int) {
		S[num[u]]--
		if S[num[u]] == 0 {
			R--
		}
		cur := R + depth
		scores[u] = cur
		if cur > best {
			best = cur
		}
		for _, v := range children[u] {
			dfsScore(v, depth+1)
		}
		if S[num[u]] == 0 {
			R++
		}
		S[num[u]]++
	}
	dfsScore(centroid, 1)
	return best, scores
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := deterministicTests()
	tests = append(tests, randomTests(40)...)

	for idx, tc := range tests {
		bestScore, scores := evaluateTree(tc.adj)
		if scores == nil {
			fmt.Fprintf(os.Stderr, "case %d: failed to evaluate tree\n", idx+1)
			os.Exit(1)
		}
		oracleOut, err := runProgram(oracle, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		oracleAns, err := parseAnswer(oracleOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid oracle output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if scores[oracleAns] != bestScore {
			fmt.Fprintf(os.Stderr, "case %d: verifier logic mismatch with oracle (score %d vs best %d)\n", idx+1, scores[oracleAns], bestScore)
			os.Exit(1)
		}
		userOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		userAns, err := parseAnswer(userOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid output: %v\n", idx+1, err)
			os.Exit(1)
		}
		if scores[userAns] != bestScore {
			fmt.Fprintf(os.Stderr, "case %d: vertex %d is not optimal (score %d, best %d)\n", idx+1, userAns, scores[userAns], bestScore)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
