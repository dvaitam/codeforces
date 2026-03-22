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

// Embedded solver for 925D

type node925 struct {
	l, r *node925
	sum  int
}

func getSum925(n *node925) int {
	if n == nil {
		return 0
	}
	return n.sum
}

func update925(root *node925, l, r, idx, val int) *node925 {
	if l == r {
		return &node925{sum: val}
	}
	mid := l + (r-l)/2
	res := &node925{}
	if root != nil {
		res.l = root.l
		res.r = root.r
	}
	if idx <= mid {
		res.l = update925(res.l, l, mid, idx, val)
	} else {
		res.r = update925(res.r, mid+1, r, idx, val)
	}
	res.sum = getSum925(res.l) + getSum925(res.r)
	return res
}

func query925(root *node925, l, r, idx int) int {
	if root == nil {
		return 0
	}
	if l == r {
		return root.sum
	}
	mid := l + (r-l)/2
	if idx <= mid {
		return query925(root.l, l, mid, idx)
	}
	return query925(root.r, mid+1, r, idx)
}

func getOnes925(root *node925, l, r int, ones *[]int) {
	if root == nil || root.sum == 0 {
		return
	}
	if l == r {
		*ones = append(*ones, l)
		return
	}
	mid := l + (r-l)/2
	getOnes925(root.l, l, mid, ones)
	getOnes925(root.r, mid+1, r, ones)
}

type state925 struct {
	u    int
	c    int
	root *node925
}

type info925 struct {
	prevU int
	prevC int
}

func solve925D(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(reader, &n, &m)

	adj := make([][]int, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	visited := make([][2]bool, n+1)
	parent := make([][2]info925, n+1)

	next := make([]int, n+2)
	prev := make([]int, n+2)
	for i := 2; i <= n; i++ {
		next[i] = i + 1
		prev[i] = i - 1
	}
	head := 2
	prev[2] = 0
	next[n] = n + 1

	removeNode := func(v int) {
		p := prev[v]
		nxt := next[v]
		if p != 0 {
			next[p] = nxt
		} else {
			head = nxt
		}
		if nxt != n+1 && nxt != 0 {
			prev[nxt] = p
		}
		prev[v] = 0
		next[v] = 0
	}

	isAdj := make([]bool, n+1)

	q := []state925{{u: 1, c: 0, root: nil}}
	visited[1][0] = true

	found := false
	for len(q) > 0 {
		curr := q[0]
		q = q[1:]

		u := curr.u
		c := curr.c
		root := curr.root

		if u == n {
			found = true
			break
		}

		for _, v := range adj[u] {
			isAdj[v] = true
		}

		newRoot := update925(root, 1, n, u, 1-c)

		for _, v := range adj[u] {
			dv := query925(root, 1, n, v)
			if dv == c {
				if !visited[v][c] {
					visited[v][c] = true
					parent[v][c] = info925{u, c}
					q = append(q, state925{v, c, newRoot})
					if c == 0 && prev[v] != 0 || head == v {
						removeNode(v)
					}
				}
			}
		}

		if c == 0 {
			var ones []int
			getOnes925(root, 1, n, &ones)
			for _, v := range ones {
				if v == u {
					continue
				}
				if !isAdj[v] {
					if !visited[v][1] {
						visited[v][1] = true
						parent[v][1] = info925{u, c}
						q = append(q, state925{v, 1, newRoot})
					}
				}
			}
		} else {
			var toRemove []int
			for v := head; v != 0 && v <= n; v = next[v] {
				if v == u {
					continue
				}
				if !isAdj[v] {
					if query925(root, 1, n, v) == 0 {
						visited[v][0] = true
						toRemove = append(toRemove, v)
						parent[v][0] = info925{u, c}
						q = append(q, state925{v, 0, newRoot})
					}
				}
			}
			for _, v := range toRemove {
				removeNode(v)
			}
		}

		for _, v := range adj[u] {
			isAdj[v] = false
		}
	}

	if !found && !visited[n][0] && !visited[n][1] {
		return "-1"
	}

	currU := n
	currC := 0
	if !visited[n][0] {
		currC = 1
	}

	path := []int{n}
	for currU != 1 || currC != 0 {
		p := parent[currU][currC]
		path = append(path, p.prevU)
		currU = p.prevU
		currC = p.prevC
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, len(path)-1)
	for i := len(path) - 1; i >= 0; i-- {
		if i < len(path)-1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", path[i])
	}
	return strings.TrimSpace(sb.String())
}

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	n := rand.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	m := rand.Intn(maxEdges + 1)
	edges := make(map[[2]int]struct{})
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for len(edges) < m {
		u := rand.Intn(n) + 1
		v := rand.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if _, ok := edges[key]; ok {
			continue
		}
		edges[key] = struct{}{}
		sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want := solve925D(string(input))

		gotRaw, err := run(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		got := strings.TrimSpace(gotRaw)

		if want == "-1" {
			if got != "-1" {
				fmt.Printf("wrong answer on test %d\n", i+1)
				fmt.Println("input:\n", string(input))
				fmt.Println("expected:\n", want)
				fmt.Println("got:\n", got)
				os.Exit(1)
			}
			continue
		}
		// For non -1 answers, compare the path length (first line)
		wantLines := strings.Split(want, "\n")
		gotLines := strings.Split(got, "\n")
		if len(wantLines) < 1 || len(gotLines) < 1 {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
		if strings.TrimSpace(wantLines[0]) != strings.TrimSpace(gotLines[0]) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
