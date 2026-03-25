package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ---------- embedded solver (from cf_t25_87_D.go) ----------

type Edge87 struct {
	u   int
	v   int
	w   int64
	idx int
}

type Adj87 struct {
	to  int
	idx int
}

type DSU87 struct {
	p  []int
	sz []int
}

func NewDSU87(n int) *DSU87 {
	p := make([]int, n+1)
	sz := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU87{p: p, sz: sz}
}

func (d *DSU87) Find(x int) int {
	for d.p[x] != x {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *DSU87) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.sz[ra] += d.sz[rb]
}

func solveD87(input string) string {
	data := []byte(input)
	pos := 0
	nextInt := func() int {
		for pos < len(data) && (data[pos] < '0' || data[pos] > '9') {
			pos++
		}
		val := 0
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			val = val*10 + int(data[pos]-'0')
			pos++
		}
		return val
	}

	n := nextInt()
	edges := make([]Edge87, n-1)
	for i := 1; i <= n-1; i++ {
		a := nextInt()
		b := nextInt()
		d := nextInt()
		edges[i-1] = Edge87{u: a, v: b, w: int64(d), idx: i}
	}

	sort.Slice(edges, func(i, j int) bool {
		if edges[i].w == edges[j].w {
			return edges[i].idx < edges[j].idx
		}
		return edges[i].w < edges[j].w
	})

	dsu := NewDSU87(n)
	ans := make([]int64, n)

	for i := 0; i < len(edges); {
		j := i + 1
		for j < len(edges) && edges[j].w == edges[i].w {
			j++
		}

		mp := make(map[int]int, 2*(j-i))
		nodeSize := make([]int64, 0, 2*(j-i))
		adj := make([][]Adj87, 0, 2*(j-i))

		getID := func(rep int) int {
			if id, ok := mp[rep]; ok {
				return id
			}
			id := len(nodeSize)
			mp[rep] = id
			nodeSize = append(nodeSize, int64(dsu.sz[rep]))
			adj = append(adj, nil)
			return id
		}

		for k := i; k < j; k++ {
			ru := dsu.Find(edges[k].u)
			rv := dsu.Find(edges[k].v)
			iu := getID(ru)
			iv := getID(rv)
			eid := edges[k].idx
			adj[iu] = append(adj[iu], Adj87{to: iv, idx: eid})
			adj[iv] = append(adj[iv], Adj87{to: iu, idx: eid})
		}

		m := len(nodeSize)
		parent := make([]int, m)
		parentEdge := make([]int, m)
		subtree := make([]int64, m)
		for t := 0; t < m; t++ {
			parent[t] = -2
		}

		stack := make([]int, 0, m)
		order := make([]int, 0, m)

		for root := 0; root < m; root++ {
			if parent[root] != -2 {
				continue
			}
			total := int64(0)
			stack = stack[:0]
			order = order[:0]
			parent[root] = -1
			stack = append(stack, root)

			for len(stack) > 0 {
				v := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				order = append(order, v)
				total += nodeSize[v]
				for _, e := range adj[v] {
					if parent[e.to] != -2 {
						continue
					}
					parent[e.to] = v
					parentEdge[e.to] = e.idx
					stack = append(stack, e.to)
				}
			}

			for _, v := range order {
				subtree[v] = nodeSize[v]
			}
			for t := len(order) - 1; t > 0; t-- {
				v := order[t]
				s := subtree[v]
				ans[parentEdge[v]] = s * (total - s)
				subtree[parent[v]] += s
			}
		}

		for k := i; k < j; k++ {
			dsu.Union(edges[k].u, edges[k].v)
		}

		i = j
	}

	maxVal := int64(0)
	cnt := 0
	for i := 1; i <= n-1; i++ {
		if ans[i] > maxVal {
			maxVal = ans[i]
			cnt = 1
		} else if ans[i] == maxVal {
			cnt++
		}
	}

	var out bytes.Buffer
	out.WriteString(strconv.FormatInt(maxVal*2, 10))
	out.WriteByte(' ')
	out.WriteString(strconv.Itoa(cnt))
	out.WriteByte('\n')

	first := true
	for i := 1; i <= n-1; i++ {
		if ans[i] == maxVal {
			if !first {
				out.WriteByte(' ')
			}
			first = false
			out.WriteString(strconv.Itoa(i))
		}
	}
	out.WriteByte('\n')

	return strings.TrimSpace(out.String())
}

// ---------- end embedded solver ----------

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(15) + 2
	type edge struct{ u, v, w int }
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		u := i
		v := rng.Intn(i-1) + 1
		w := rng.Intn(100) + 1
		edges = append(edges, edge{u, v, w})
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	// suppress unused import
	_ = io.Discard

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp := solveD87(input)
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\n got:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
