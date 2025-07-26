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

func runCandidate(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

type edge struct{ to, val, next int }

func expected(n, s int, edgesList [][3]int) int64 {
	head := make([]int, n+1)
	edges := make([]edge, 2*(n+1))
	ec := 0
	for _, e := range edgesList {
		u, v, w := e[0], e[1], e[2]
		ec++
		edges[ec] = edge{v, w, head[u]}
		head[u] = ec
		ec++
		edges[ec] = edge{u, w, head[v]}
		head[v] = ec
	}
	dis := make([]int64, n+1)
	fa := make([]int, n+1)
	type pair struct{ x, p int }
	var farthest func(int) int
	farthest = func(st int) int {
		for i := 1; i <= n; i++ {
			dis[i] = 0
			fa[i] = 0
		}
		stack := []pair{{st, 0}}
		mx := st
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			x, p := cur.x, cur.p
			if dis[x] > dis[mx] {
				mx = x
			}
			for e := head[x]; e != 0; e = edges[e].next {
				y := edges[e].to
				if y == p {
					continue
				}
				dis[y] = dis[x] + int64(edges[e].val)
				fa[y] = x
				stack = append(stack, pair{y, x})
			}
		}
		return mx
	}
	u := farthest(1)
	v := farthest(u)
	path := []int{}
	on := make([]bool, n+1)
	for x := v; x != 0; x = fa[x] {
		path = append(path, x)
		on[x] = true
	}
	m := len(path)
	disU := make([]int64, m)
	for i, x := range path {
		disU[i] = dis[x]
	}
	mxd := make([]int64, m)
	type info struct {
		x, p int
		d    int64
	}
	for i, x := range path {
		var b int64
		for e := head[x]; e != 0; e = edges[e].next {
			y := edges[e].to
			if on[y] {
				continue
			}
			st := []info{{y, x, int64(edges[e].val)}}
			for len(st) > 0 {
				cur := st[len(st)-1]
				st = st[:len(st)-1]
				if cur.d > b {
					b = cur.d
				}
				for ee := head[cur.x]; ee != 0; ee = edges[ee].next {
					yy := edges[ee].to
					if yy == cur.p {
						continue
					}
					st = append(st, info{yy, cur.x, cur.d + int64(edges[ee].val)})
				}
			}
		}
		mxd[i] = b
	}
	best := int64(1<<63 - 1)
	dq := []int{}
	r := 0
	for l := 0; l < m; l++ {
		if len(dq) > 0 && dq[0] < l {
			dq = dq[1:]
		}
		for r < m && r-l < s {
			for len(dq) > 0 && mxd[r] >= mxd[dq[len(dq)-1]] {
				dq = dq[:len(dq)-1]
			}
			dq = append(dq, r)
			r++
		}
		cand := max(max(disU[0]-disU[l], disU[r-1]), mxd[dq[0]])
		if cand < best {
			best = cand
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n, s  int
		edges [][3]int
	}
	var cases []test
	cases = append(cases, test{1, 1, [][3]int{}})
	cases = append(cases, test{2, 1, [][3]int{{1, 2, 5}}})
	for i := 0; i < 98; i++ {
		n := rng.Intn(20) + 1
		s := rng.Intn(n) + 1
		edges := make([][3]int, n-1)
		for j := 1; j < n; j++ {
			u := j + 1
			v := rng.Intn(j) + 1
			w := rng.Intn(10) + 1
			edges[j-1] = [3]int{u, v, w}
		}
		cases = append(cases, test{n, s, edges})
	}

	for idx, tc := range cases {
		input := fmt.Sprintf("%d %d\n", tc.n, tc.s)
		for _, e := range tc.edges {
			input += fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2])
		}
		want := fmt.Sprintf("%d", expected(tc.n, tc.s, tc.edges))
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
