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

type Edge struct{ u, v, w int }

type subset struct{ edges []int }

func genGraph(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	m := rng.Intn(3) + n - 1
	a := rng.Intn(5) + 1
	b := a + rng.Intn(5) + 1
	edges := make([]Edge, m)
	used := make([]bool, m)
	// ensure connectivity via tree first
	for i := 1; i < n; i++ {
		u := i
		v := rng.Intn(i) + 1
		w := []int{a, b}[rng.Intn(2)]
		edges[i-1] = Edge{u, v, w}
		used[i-1] = true
	}
	for i := n - 1; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for u == v {
			v = rng.Intn(n) + 1
		}
		w := []int{a, b}[rng.Intn(2)]
		edges[i] = Edge{u, v, w}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, a, b)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	// compute expected output via brute force
	bestSum := int(1 << 30)
	type tree struct {
		sum  int
		dist []int
	}
	var bestTrees []tree
	idxs := make([]int, m)
	for i := 0; i < m; i++ {
		idxs[i] = i
	}
	subsetSz := n - 1
	var choose func(int, int, []int)
	choose = func(pos, chosen int, cur []int) {
		if chosen == subsetSz {
			// check if edges in cur form tree
			adj := make([][]Edge, n)
			for _, id := range cur {
				e := edges[id]
				adj[e.u-1] = append(adj[e.u-1], Edge{e.v - 1, 0, e.w})
				adj[e.v-1] = append(adj[e.v-1], Edge{e.u - 1, 0, e.w})
			}
			// BFS connectivity
			vis := make([]bool, n)
			stack := []int{0}
			vis[0] = true
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				for _, ed := range adj[u] {
					if !vis[ed.u] {
						vis[ed.u] = true
						stack = append(stack, ed.u)
					}
				}
			}
			all := true
			for _, v := range vis {
				if !v {
					all = false
					break
				}
			}
			if !all {
				return
			}
			sum := 0
			for _, id := range cur {
				sum += edges[id].w
			}
			if sum < bestSum {
				bestSum = sum
				bestTrees = nil
			}
			if sum == bestSum {
				// compute dist from 1
				dist := make([]int, n)
				for i := 0; i < n; i++ {
					dist[i] = 1 << 30
				}
				dist[0] = 0
				q := []int{0}
				for len(q) > 0 {
					u := q[0]
					q = q[1:]
					for _, ed := range adj[u] {
						if dist[ed.u] > dist[u]+ed.w {
							dist[ed.u] = dist[u] + ed.w
							q = append(q, ed.u)
						}
					}
				}
				bestTrees = append(bestTrees, tree{sum, dist})
			}
			return
		}
		if pos == m {
			return
		}
		// choose edge pos
		cur2 := append(cur, idxs[pos])
		choose(pos+1, chosen+1, cur2)
		choose(pos+1, chosen, cur)
	}
	choose(0, 0, nil)
	outVals := make([]int, n)
	for i := 0; i < n; i++ {
		best := 1 << 30
		for _, tr := range bestTrees {
			if tr.dist[i] < best {
				best = tr.dist[i]
			}
		}
		outVals[i] = best
	}
	var out strings.Builder
	for i, v := range outVals {
		if i > 0 {
			out.WriteByte(' ')
		}
		out.WriteString(fmt.Sprintf("%d", v))
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in, expect := genGraph(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
