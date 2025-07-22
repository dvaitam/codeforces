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

type query struct {
	typ int
	a   int
	b   int
}

func run(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func recomputeDepth(children [][]int, parent []int, depth []int) {
	queue := []int{1}
	depth[1] = 0
	for i := 0; i < len(queue); i++ {
		u := queue[i]
		for _, v := range children[u] {
			depth[v] = depth[u] + 1
			queue = append(queue, v)
		}
	}
}

func lca(u, v int, parent, depth []int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	for depth[u] > depth[v] {
		u = parent[u]
	}
	for u != v {
		u = parent[u]
		v = parent[v]
	}
	return u
}

func solveE(n int, children [][]int, parent []int, qs []query) []int {
	depth := make([]int, n+1)
	recomputeDepth(children, parent, depth)
	var res []int
	for _, q := range qs {
		switch q.typ {
		case 1:
			u := q.a
			v := q.b
			w := lca(u, v, parent, depth)
			dist := depth[u] + depth[v] - 2*depth[w]
			res = append(res, dist)
		case 2:
			v := q.a
			h := q.b
			p := v
			for i := 0; i < h; i++ {
				p = parent[p]
			}
			old := parent[v]
			lst := children[old]
			for i, x := range lst {
				if x == v {
					children[old] = append(lst[:i], lst[i+1:]...)
					break
				}
			}
			children[p] = append(children[p], v)
			parent[v] = p
			recomputeDepth(children, parent, depth)
		case 3:
			k := q.a
			ans := -1
			var dfs func(int)
			dfs = func(u int) {
				if depth[u] == k {
					ans = u
				}
				for _, v := range children[u] {
					dfs(v)
				}
			}
			dfs(1)
			res = append(res, ans)
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(4) + 2
	children := make([][]int, n+1)
	parent := make([]int, n+1)
	for v := 2; v <= n; v++ {
		p := rng.Intn(v-1) + 1
		children[p] = append(children[p], v)
		parent[v] = p
	}
	parent[1] = 0
	m := rng.Intn(8) + 1
	qs := make([]query, 0, m)
	depth := make([]int, n+1)
	recomputeDepth(children, parent, depth)
	for i := 0; i < m; i++ {
		typ := rng.Intn(3) + 1
		switch typ {
		case 1:
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			qs = append(qs, query{1, u, v})
		case 2:
			// ensure valid vertex with depth>=2
			cand := []int{}
			for v := 2; v <= n; v++ {
				if depth[v] >= 2 {
					cand = append(cand, v)
				}
			}
			if len(cand) == 0 {
				typ = 1
				u := rng.Intn(n) + 1
				v := rng.Intn(n) + 1
				qs = append(qs, query{1, u, v})
				continue
			}
			v := cand[rng.Intn(len(cand))]
			h := rng.Intn(depth[v]-1) + 2
			qs = append(qs, query{2, v, h})
			// apply to maintain validity
			p := v
			for j := 0; j < h; j++ {
				p = parent[p]
			}
			old := parent[v]
			lst := children[old]
			for idx, x := range lst {
				if x == v {
					children[old] = append(lst[:idx], lst[idx+1:]...)
					break
				}
			}
			children[p] = append(children[p], v)
			parent[v] = p
			recomputeDepth(children, parent, depth)
		case 3:
			// pick random depth value that exists
			depthMap := map[int]struct{}{}
			for x := 1; x <= n; x++ {
				depthMap[depth[x]] = struct{}{}
			}
			vals := make([]int, 0, len(depthMap))
			for d := range depthMap {
				vals = append(vals, d)
			}
			k := vals[rng.Intn(len(vals))]
			qs = append(qs, query{3, k, 0})
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(qs)))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d", len(children[i])))
		for _, ch := range children[i] {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(ch))
		}
		sb.WriteByte('\n')
	}
	for _, q := range qs {
		if q.typ == 1 {
			sb.WriteString(fmt.Sprintf("1 %d %d\n", q.a, q.b))
		} else if q.typ == 2 {
			sb.WriteString(fmt.Sprintf("2 %d %d\n", q.a, q.b))
		} else {
			sb.WriteString(fmt.Sprintf("3 %d\n", q.a))
		}
	}
	// compute expected using solver on a fresh copy
	childrenCopy := make([][]int, n+1)
	parentCopy := make([]int, n+1)
	for i := 1; i <= n; i++ {
		parentCopy[i] = parent[i]
		if len(children[i]) > 0 {
			childrenCopy[i] = append([]int(nil), children[i]...)
		}
	}
	expected := solveE(n, childrenCopy, parentCopy, qs)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		outTok := strings.Fields(out)
		if len(outTok) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d outputs got %d\ninput:\n%soutput:\n%s", i+1, len(exp), len(outTok), input, out)
			os.Exit(1)
		}
		for j, tok := range outTok {
			got, err := strconv.Atoi(tok)
			if err != nil || got != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at output %d: expected %d got %s\ninput:\n%s", i+1, j+1, exp[j], tok, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
