package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runProg(bin, input string) (string, error) {
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

func expect(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	readInt := func() int {
		in.Scan()
		v := 0
		fmt.Sscan(in.Text(), &v)
		return v
	}
	n := readInt()
	type edge struct{ u, v int }
	edges := make([]edge, n-1)
	adj := make([][]int, n)
	for i := 0; i < n-1; i++ {
		u := readInt() - 1
		v := readInt() - 1
		edges[i] = edge{u, v}
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	parent := make([]int, n)
	depth := make([]int, n)
	parent[0] = -1
	queue := []int{0}
	for idx := 0; idx < len(queue); idx++ {
		u := queue[idx]
		for _, v := range adj[u] {
			if v == parent[u] {
				continue
			}
			parent[v] = u
			depth[v] = depth[u] + 1
			queue = append(queue, v)
		}
	}
	value := make([]int, n)
	for i := range value {
		value[i] = 1
	}
	m := readInt()
	type cons struct{ u, v, w int }
	consList := make([]cons, m)
	for i := 0; i < m; i++ {
		u := readInt() - 1
		v := readInt() - 1
		w := readInt()
		consList[i] = cons{u, v, w}
		uu, vv := u, v
		if depth[uu] < depth[vv] {
			uu, vv = vv, uu
		}
		for depth[uu] > depth[vv] {
			if value[uu] < w {
				value[uu] = w
			}
			uu = parent[uu]
		}
		for uu != vv {
			if value[uu] < w {
				value[uu] = w
			}
			if value[vv] < w {
				value[vv] = w
			}
			uu = parent[uu]
			vv = parent[vv]
		}
	}
	ok := true
	for _, c := range consList {
		u, v, w := c.u, c.v, c.w
		uu, vv := u, v
		if depth[uu] < depth[vv] {
			uu, vv = vv, uu
		}
		found := false
		for depth[uu] > depth[vv] {
			if value[uu] == w {
				found = true
				break
			}
			uu = parent[uu]
		}
		for !found && uu != vv {
			if value[uu] == w || value[vv] == w {
				found = true
				break
			}
			uu = parent[uu]
			vv = parent[vv]
		}
		if !found {
			ok = false
			break
		}
	}
	if !ok {
		return "-1"
	}
	ans := make([]int, n-1)
	for i, e := range edges {
		if parent[e.u] == e.v {
			ans[i] = value[e.u]
		} else {
			ans[i] = value[e.v]
		}
	}
	var sb strings.Builder
	for i, w := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", w))
	}
	return sb.String()
}

type testCase struct{ input string }

var tests = []testCase{
	{"2\n1 2\n0\n"},
	{"2\n1 2\n1\n1 2 2\n"},
	{"3\n1 2\n2 3\n1\n1 3 3\n"},
	{"3\n1 2\n2 3\n2\n1 3 3\n1 2 2\n"},
	{"5\n1 2\n2 3\n3 4\n4 5\n3\n1 3 2\n2 4 3\n5 3 1\n"},
	{"4\n1 2\n2 3\n2 4\n1\n3 4 2\n"},
	{"4\n1 2\n1 3\n3 4\n1\n2 4 2\n"},
	{"4\n1 2\n2 3\n1 4\n1\n1 3 2\n"},
	{"3\n1 2\n1 3\n2\n1 2 2\n2 3 2\n"},
	{"3\n1 2\n2 3\n3\n2 3 2\n1 3 2\n1 3 1\n"},
	{"5\n1 2\n1 3\n1 4\n1 5\n2\n2 3 2\n4 5 3\n"},
	{"3\n1 2\n1 3\n0\n"},
	{"2\n1 2\n0\n"},
	{"2\n1 2\n1\n1 2 3\n"},
	{"4\n1 2\n1 3\n1 4\n2\n2 3 1\n3 4 2\n"},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range tests {
		exp := expect(tc.input)
		got, err := runProg(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
