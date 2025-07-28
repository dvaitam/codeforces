package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func solveOne(n, m int, val []int64, edges [][2]int) (int, int64) {
	g := make([][]int, n)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
	}

	index := make([]int, n)
	low := make([]int, n)
	onStack := make([]bool, n)
	stack := make([]int, 0, n)
	compID := make([]int, n)
	compCnt := 0
	idx := 0
	compSize := []int{}
	compSum := []int64{}

	var dfs func(int)
	dfs = func(v int) {
		idx++
		index[v] = idx
		low[v] = idx
		stack = append(stack, v)
		onStack[v] = true
		for _, to := range g[v] {
			if index[to] == 0 {
				dfs(to)
				if low[to] < low[v] {
					low[v] = low[to]
				}
			} else if onStack[to] && index[to] < low[v] {
				low[v] = index[to]
			}
		}
		if low[v] == index[v] {
			size := 0
			sum := int64(0)
			for {
				w := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				onStack[w] = false
				compID[w] = compCnt
				size++
				sum += val[w]
				if w == v {
					break
				}
			}
			compSize = append(compSize, size)
			compSum = append(compSum, sum)
			compCnt++
		}
	}
	for i := 0; i < n; i++ {
		if index[i] == 0 {
			dfs(i)
		}
	}

	dag := make([]map[int]struct{}, compCnt)
	for i := 0; i < compCnt; i++ {
		dag[i] = make(map[int]struct{})
	}
	for v := 0; v < n; v++ {
		cv := compID[v]
		for _, u := range g[v] {
			cu := compID[u]
			if cv != cu {
				dag[cv][cu] = struct{}{}
			}
		}
	}

	adj := make([][]int, compCnt)
	indeg := make([]int, compCnt)
	for i := 0; i < compCnt; i++ {
		for nb := range dag[i] {
			adj[i] = append(adj[i], nb)
			indeg[nb]++
		}
	}

	order := make([]int, 0, compCnt)
	q := make([]int, 0, compCnt)
	for i := 0; i < compCnt; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for qi := 0; qi < len(q); qi++ {
		v := q[qi]
		order = append(order, v)
		for _, u := range adj[v] {
			indeg[u]--
			if indeg[u] == 0 {
				q = append(q, u)
			}
		}
	}

	dpLen := make([]int, compCnt)
	dpSum := make([]int64, compCnt)
	for i := compCnt - 1; i >= 0; i-- {
		v := order[i]
		dpLen[v] = compSize[v]
		dpSum[v] = compSum[v]
		for _, u := range adj[v] {
			candLen := compSize[v] + dpLen[u]
			candSum := compSum[v] + dpSum[u]
			if candLen > dpLen[v] {
				dpLen[v] = candLen
				dpSum[v] = candSum
			} else if candLen == dpLen[v] && candSum < dpSum[v] {
				dpSum[v] = candSum
			}
		}
	}

	bestLen := 0
	bestSum := int64(0)
	first := true
	for i := 0; i < compCnt; i++ {
		if first || dpLen[i] > bestLen || (dpLen[i] == bestLen && dpSum[i] < bestSum) {
			bestLen = dpLen[i]
			bestSum = dpSum[i]
			first = false
		}
	}
	return bestLen, bestSum
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(3) + 1
	var in strings.Builder
	var out strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		m := rng.Intn(n*n + 1)
		vals := make([]int64, n)
		for j := 0; j < n; j++ {
			vals[j] = int64(rng.Intn(20))
		}
		edges := make([][2]int, m)
		for j := 0; j < m; j++ {
			edges[j][0] = rng.Intn(n)
			edges[j][1] = rng.Intn(n)
		}
		in.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for j := 0; j < n; j++ {
			if j > 0 {
				in.WriteByte(' ')
			}
			in.WriteString(fmt.Sprintf("%d", vals[j]))
		}
		in.WriteByte('\n')
		for _, e := range edges {
			in.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}
		l, s := solveOne(n, m, vals, edges)
		out.WriteString(fmt.Sprintf("%d %d\n", l, s))
	}
	return testCase{input: in.String(), expected: out.String()}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(tc.expected)
	if got != exp {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	// simple case with 1 node
	in := "1\n1 0\n5\n"
	l, s := solveOne(1, 0, []int64{5}, nil)
	out := fmt.Sprintf("%d %d\n", l, s)
	cases := []testCase{{input: in, expected: out}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
