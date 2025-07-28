package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type segment struct {
	x, l, r int
}

type Node struct {
	key, prio   int
	left, right *Node
}

func split(root *Node, key int) (l, r *Node) {
	if root == nil {
		return nil, nil
	}
	if root.key < key {
		l1, r1 := split(root.right, key)
		root.right = l1
		return root, r1
	}
	l1, r1 := split(root.left, key)
	root.left = r1
	return l1, root
}

func merge(l, r *Node) *Node {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.prio > r.prio {
		l.right = merge(l.right, r)
		return l
	}
	r.left = merge(l, r.left)
	return r
}

func insert(root *Node, node *Node) *Node {
	if root == nil {
		return node
	}
	if node.prio > root.prio {
		l, r := split(root, node.key)
		node.left, node.right = l, r
		return node
	}
	if node.key < root.key {
		root.left = insert(root.left, node)
	} else {
		root.right = insert(root.right, node)
	}
	return root
}

func erase(root *Node, key int) *Node {
	if root == nil {
		return nil
	}
	if key == root.key {
		return merge(root.left, root.right)
	}
	if key < root.key {
		root.left = erase(root.left, key)
	} else {
		root.right = erase(root.right, key)
	}
	return root
}

func predecessor(root *Node, key int) int {
	res := -1000000000
	for root != nil {
		if root.key < key {
			if root.key > res {
				res = root.key
			}
			root = root.right
		} else {
			root = root.left
		}
	}
	return res
}

func successor(root *Node, key int) int {
	res := 1000000000
	for root != nil {
		if root.key > key {
			if root.key < res {
				res = root.key
			}
			root = root.left
		} else {
			root = root.right
		}
	}
	return res
}

func solveD(n int, segs []segment) string {
	rand.Seed(1)
	cover := make([]int, n)
	type evt struct{ pos, x, t int }
	endpoints := make([]evt, 0, len(segs)*2)
	for _, s := range segs {
		x := s.x - 1
		l := s.l - 1
		r := s.r
		endpoints = append(endpoints, evt{l, x, -1})
		endpoints = append(endpoints, evt{r, x, 1})
	}
	sort.Slice(endpoints, func(i, j int) bool { return endpoints[i].pos < endpoints[j].pos })
	var root *Node
	root = insert(root, &Node{key: -1, prio: rand.Int()})
	root = insert(root, &Node{key: n, prio: rand.Int()})
	edges := make([][2]int, 0, len(segs)*4)
	add := make([]int, 0)
	del := make([]int, 0)
	for l := 0; l < len(endpoints); {
		r := l
		for r < len(endpoints) && endpoints[r].pos == endpoints[l].pos {
			r++
		}
		add = add[:0]
		del = del[:0]
		for i := l; i < r; i++ {
			e := endpoints[i]
			if e.t == -1 {
				if cover[e.x] == 0 {
					add = append(add, e.x)
				}
				cover[e.x]++
			} else {
				cover[e.x]--
				if cover[e.x] == 0 {
					del = append(del, e.x)
				}
			}
		}
		for _, x := range del {
			root = erase(root, x)
		}
		for _, x := range add {
			root = insert(root, &Node{key: x, prio: rand.Int()})
			p := predecessor(root, x)
			s := successor(root, x)
			edges = append(edges, [2]int{p, x})
			edges = append(edges, [2]int{x, s})
		}
		l = r
	}
	sort.Slice(edges, func(i, j int) bool {
		if edges[i][0] != edges[j][0] {
			return edges[i][0] < edges[j][0]
		}
		return edges[i][1] < edges[j][1]
	})
	dp := make([]int, n+2)
	g := make([]int, n+2)
	for _, e := range edges {
		u, v := e[0], e[1]
		ui, vi := u+1, v+1
		if dp[ui]+1 > dp[vi] {
			dp[vi] = dp[ui] + 1
			g[vi] = ui
		}
	}
	res := n + 1 - dp[n+1]
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", res))
	list := make([]int, 0, res)
	for i := n + 1; i > 0; {
		prev := g[i]
		for j := prev + 1; j < i; j++ {
			list = append(list, j)
		}
		i = prev
	}
	for _, v := range list {
		sb.WriteString(fmt.Sprintf("%d\n", v))
	}
	return strings.TrimSpace(sb.String())
}

func generateD(rng *rand.Rand) (int, []segment) {
	n := rng.Intn(8) + 1
	m := rng.Intn(10)
	segs := make([]segment, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n) + 1
		l := rng.Intn(50) + 1
		r := l + rng.Intn(50)
		segs[i] = segment{x, l, r}
	}
	return n, segs
}

func runCase(bin string, n int, segs []segment) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(segs)))
	for _, s := range segs {
		fmt.Fprintf(&sb, "%d %d %d\n", s.x, s.l, s.r)
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]struct {
		n    int
		segs []segment
	}, 0, 102)
	// small fixed cases
	cases = append(cases, struct {
		n    int
		segs []segment
	}{1, nil})
	cases = append(cases, struct {
		n    int
		segs []segment
	}{2, []segment{{1, 1, 1}}})
	for i := 0; i < 100; i++ {
		n, segs := generateD(rng)
		cases = append(cases, struct {
			n    int
			segs []segment
		}{n, segs})
	}
	for i, tc := range cases {
		expect := solveD(tc.n, tc.segs)
		out, err := runCase(bin, tc.n, tc.segs)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\n\nGot:\n%s\n", i+1, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
