package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxLog = 21 // since n <= 1e6, 2^20 > 1e6

type fastScanner struct {
	r *bufio.Reader
}

func newScanner() *fastScanner {
	return &fastScanner{r: bufio.NewReader(os.Stdin)}
}

func (fs *fastScanner) nextInt() int {
	sign, val := 1, 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		c, err = fs.r.ReadByte()
		if err != nil {
			return 0
		}
	}
	if c == '-' {
		sign = -1
		c, _ = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := newScanner()
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	t := in.nextInt()
	for ; t > 0; t-- {
		n := in.nextInt()
		st := in.nextInt()

		w := make([]int, n+1)
		for i := 1; i <= n; i++ {
			w[i] = in.nextInt()
		}

		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			u := in.nextInt()
			v := in.nextInt()
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}

		depth := make([]int, n+1)
		up := make([][]int, maxLog)
		for i := range up {
			up[i] = make([]int, n+1)
		}

		// DFS from root 1 to fill depth and up table
		stack := []int{1}
		parent := make([]int, n+1)
		parent[1] = 0
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			up[0][v] = parent[v]
			for _, to := range adj[v] {
				if to == parent[v] {
					continue
				}
				parent[to] = v
				depth[to] = depth[v] + 1
				stack = append(stack, to)
			}
		}

		for k := 1; k < maxLog; k++ {
			for v := 1; v <= n; v++ {
				up[k][v] = up[k-1][up[k-1][v]]
			}
		}

		lca := func(a, b int) int {
			if depth[a] < depth[b] {
				a, b = b, a
			}
			diff := depth[a] - depth[b]
			for k := 0; diff > 0; k++ {
				if diff&1 == 1 {
					a = up[k][a]
				}
				diff >>= 1
			}
			if a == b {
				return a
			}
			for k := maxLog - 1; k >= 0; k-- {
				if up[k][a] != up[k][b] {
					a = up[k][a]
					b = up[k][b]
				}
			}
			return up[0][a]
		}

		// DFS from st to compute cumulative sums and minimal prefixes along the path from st
		cum := make([]int, n+1)     // sum of weights from st to node (inclusive)
		minPref := make([]int, n+1) // minimal prefix sum along path
		parST := make([]int, n+1)   // parent in the st-rooted tree

		type nodeInfo struct {
			v, p int
		}
		stStack := []nodeInfo{{st, 0}}
		cum[st] = w[st]
		minPref[st] = w[st]

		for len(stStack) > 0 {
			cur := stStack[len(stStack)-1]
			stStack = stStack[:len(stStack)-1]
			v := cur.v
			p := cur.p
			parST[v] = p
			for _, to := range adj[v] {
				if to == p {
					continue
				}
				cum[to] = cum[v] + w[to]
				minPref[to] = minInt(minPref[v], cum[to])
				stStack = append(stStack, nodeInfo{to, v})
			}
		}

		ans := 0
		depthSt := depth[st]

		for v := 1; v <= n; v++ {
			if v == st {
				continue
			}
			if minPref[v] < 0 {
				continue // would die before reaching v
			}
			l := lca(st, v)
			h := 2*depth[l] - depthSt // remaining margin after reaching l
			if h <= 0 {
				continue
			}
			extra := h
			if extra%2 == 0 {
				extra--
			}
			if extra <= 0 {
				continue
			}
			dist := depthSt + depth[v] - 2*depth[l]
			curSum := cum[v] // sum of weights along the path from st to v

			wp := w[parST[v]]
			wv := w[v]
			life := extra

			if wp == -1 && wv == 1 {
				if curSum == 0 {
					life = 0
				}
			} else if wp == -1 && wv == -1 {
				if curSum < life {
					life = curSum
				}
			}

			if life < 0 {
				life = 0
			}
			if life > extra {
				life = extra
			}

			cand := dist + life
			if cand > ans {
				ans = cand
			}
		}

		fmt.Fprintln(out, ans)
	}
}
