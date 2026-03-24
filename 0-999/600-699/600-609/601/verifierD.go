package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const solverSrc = `package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	children [26]int32
	size     int32
}

var nodes []Node

func merge(a, b int32) int32 {
	if a == 0 {
		return b
	}
	if b == 0 {
		return a
	}
	nodes[a].size = 1
	for i := 0; i < 26; i++ {
		if nodes[a].children[i] == 0 {
			nodes[a].children[i] = nodes[b].children[i]
		} else if nodes[b].children[i] != 0 {
			nodes[a].children[i] = merge(nodes[a].children[i], nodes[b].children[i])
		}
		if nodes[a].children[i] != 0 {
			nodes[a].size += nodes[nodes[a].children[i]].size
		}
	}
	return a
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)

	var n int
	fmt.Fscan(in, &n)

	if n == 0 {
		return
	}

	c := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &c[i])
	}

	var s string
	fmt.Fscan(in, &s)

	head := make([]int32, n+1)
	next := make([]int32, 2*n)
	to := make([]int32, 2*n)
	var edgeCnt int32

	addEdge := func(u, v int32) {
		edgeCnt++
		next[edgeCnt] = head[u]
		to[edgeCnt] = v
		head[u] = edgeCnt
	}

	for i := 1; i < n; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		addEdge(u, v)
		addEdge(v, u)
	}

	nodes = make([]Node, n+1)
	dif := make([]int32, n+1)

	order := make([]int32, 0, n)
	q := make([]int32, 0, n)
	parent := make([]int32, n+1)

	q = append(q, 1)
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		order = append(order, u)
		for e := head[u]; e != 0; e = next[e] {
			v := to[e]
			if v != parent[u] {
				parent[v] = u
				q = append(q, v)
			}
		}
	}

	for i := len(order) - 1; i >= 0; i-- {
		u := order[i]
		nodes[u].size = 1
		for e := head[u]; e != 0; e = next[e] {
			v := to[e]
			if v != parent[u] {
				charIdx := s[v-1] - 'a'
				nodes[u].children[charIdx] = merge(nodes[u].children[charIdx], v)
			}
		}
		nodes[u].size = 1
		for j := 0; j < 26; j++ {
			if nodes[u].children[j] != 0 {
				nodes[u].size += nodes[nodes[u].children[j]].size
			}
		}
		dif[u] = nodes[u].size
	}

	var maxVal int64 = -1
	var count int

	for i := 1; i <= n; i++ {
		val := int64(dif[i]) + c[i]
		if val > maxVal {
			maxVal = val
			count = 1
		} else if val == maxVal {
			count++
		}
	}

	fmt.Printf("%d\n%d\n", maxVal, count)
}
`

func buildSolver() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verD601")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	src := filepath.Join(dir, "solver.go")
	if err := os.WriteFile(src, []byte(solverSrc), 0644); err != nil {
		cleanup()
		return "", nil, err
	}
	bin := filepath.Join(dir, "solver")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("build solver: %v\n%s", err, out)
	}
	return bin, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(63))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(6) + 1
		vals := make([]int64, n)
		for j := 0; j < n; j++ {
			vals[j] = int64(rng.Intn(10))
		}
		letters := make([]byte, n)
		for j := 0; j < n; j++ {
			letters[j] = byte('a' + rng.Intn(3))
		}
		edges := make([][2]int, n-1)
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			edges[v-2] = [2]int{p, v}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j, v := range vals {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		sb.WriteString(string(letters))
		sb.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, cleanup, err := buildSolver()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	cases := genCases()
	for i, c := range cases {
		expected, err := runBinary(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d ref failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		got, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
