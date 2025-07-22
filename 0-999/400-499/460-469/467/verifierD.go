package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type weight struct {
	rCount int
	length int
}

func (w weight) less(o weight) bool {
	if w.rCount != o.rCount {
		return w.rCount < o.rCount
	}
	return w.length < o.length
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveD(r *bufio.Reader) string {
	var m int
	if _, err := fmt.Fscan(r, &m); err != nil {
		return ""
	}
	essay := make([]string, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &essay[i])
		essay[i] = strings.ToLower(essay[i])
	}
	var n int
	fmt.Fscan(r, &n)
	id := make(map[string]int, m+2*n)
	var words []string
	addWord := func(s string) int {
		if v, ok := id[s]; ok {
			return v
		}
		idx := len(words)
		id[s] = idx
		words = append(words, s)
		return idx
	}
	for _, w := range essay {
		addWord(w)
	}
	edgePairs := make([][2]int, n)
	for i := 0; i < n; i++ {
		var x, y string
		fmt.Fscan(r, &x, &y)
		x = strings.ToLower(x)
		y = strings.ToLower(y)
		u := addWord(x)
		v := addWord(y)
		edgePairs[i] = [2]int{u, v}
	}
	N := len(words)
	edges := make([][]int, N)
	rev := make([][]int, N)
	for _, p := range edgePairs {
		u, v := p[0], p[1]
		edges[u] = append(edges[u], v)
		rev[v] = append(rev[v], u)
	}
	nodeW := make([]weight, N)
	for i, s := range words {
		cnt := 0
		for _, ch := range s {
			if ch == 'r' {
				cnt++
			}
		}
		nodeW[i] = weight{rCount: cnt, length: len(s)}
	}
	visited := make([]bool, N)
	order := make([]int, 0, N)
	type frame struct {
		v, idx int
		pre    bool
	}
	for v := 0; v < N; v++ {
		if visited[v] {
			continue
		}
		stack := []frame{{v, 0, false}}
		for len(stack) > 0 {
			fr := &stack[len(stack)-1]
			if !fr.pre {
				visited[fr.v] = true
				fr.pre = true
			}
			if fr.idx < len(edges[fr.v]) {
				u := edges[fr.v][fr.idx]
				fr.idx++
				if !visited[u] {
					stack = append(stack, frame{u, 0, false})
				}
			} else {
				order = append(order, fr.v)
				stack = stack[:len(stack)-1]
			}
		}
	}
	comp := make([]int, N)
	for i := range comp {
		comp[i] = -1
	}
	cid := 0
	for i := len(order) - 1; i >= 0; i-- {
		v := order[i]
		if comp[v] != -1 {
			continue
		}
		stk := []int{v}
		for len(stk) > 0 {
			u := stk[len(stk)-1]
			stk = stk[:len(stk)-1]
			if comp[u] != -1 {
				continue
			}
			comp[u] = cid
			for _, w := range rev[u] {
				if comp[w] == -1 {
					stk = append(stk, w)
				}
			}
		}
		cid++
	}
	C := cid
	inf := 1_000_000_000
	compW := make([]weight, C)
	for i := 0; i < C; i++ {
		compW[i] = weight{rCount: inf, length: inf}
	}
	for i := 0; i < N; i++ {
		c := comp[i]
		if nodeW[i].less(compW[c]) {
			compW[c] = nodeW[i]
		}
	}
	cdg := make([][]int, C)
	indeg := make([]int, C)
	for u := 0; u < N; u++ {
		cu := comp[u]
		for _, v := range edges[u] {
			cv := comp[v]
			if cu != cv {
				cdg[cu] = append(cdg[cu], cv)
			}
		}
	}
	for u := 0; u < C; u++ {
		for _, v := range cdg[u] {
			indeg[v]++
		}
	}
	q := make([]int, 0, C)
	for i := 0; i < C; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	topo := make([]int, 0, C)
	for i := 0; i < len(q); i++ {
		u := q[i]
		topo = append(topo, u)
		for _, v := range cdg[u] {
			indeg[v]--
			if indeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	dp := make([]weight, C)
	copy(dp, compW)
	for i := len(topo) - 1; i >= 0; i-- {
		u := topo[i]
		for _, v := range cdg[u] {
			if dp[v].less(dp[u]) {
				dp[u] = dp[v]
			}
		}
	}
	totalR, totalLen := 0, 0
	for _, w := range essay {
		c := comp[id[w]]
		totalR += dp[c].rCount
		totalLen += dp[c].length
	}
	return fmt.Sprintf("%d %d", totalR, totalLen)
}

func randomWord(rng *rand.Rand) string {
	l := rng.Intn(5) + 1
	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func generateCaseD(rng *rand.Rand) string {
	m := rng.Intn(5) + 1
	wordSet := make(map[string]struct{})
	words := make([]string, 0, m+10)
	for len(words) < m+10 {
		w := randomWord(rng)
		if _, ok := wordSet[w]; !ok {
			wordSet[w] = struct{}{}
			words = append(words, w)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", m)
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(words[i])
	}
	sb.WriteByte('\n')
	n := rng.Intn(5) + 1
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		x := words[rng.Intn(len(words))]
		y := words[rng.Intn(len(words))]
		fmt.Fprintf(&sb, "%s %s\n", x, y)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		expect := solveD(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
