package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

type Tree struct {
	parent   []int
	children [][]int
	value    []int64
}

func newTree(n int) *Tree {
	return &Tree{
		parent:   make([]int, n+1),
		children: make([][]int, n+1),
		value:    make([]int64, n+1),
	}
}

func (t *Tree) orient(root int, edges [][2]int) {
	g := make([][]int, len(t.parent))
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	stack := []int{root}
	t.parent[root] = 0
	for len(stack) > 0 {
		v := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		for _, to := range g[v] {
			if to == t.parent[v] {
				continue
			}
			t.parent[to] = v
			t.children[v] = append(t.children[v], to)
			stack = append(stack, to)
		}
	}
}

func (t *Tree) subtreeSum(v int) (int64, int) {
	sum := t.value[v]
	size := 1
	for _, c := range t.children[v] {
		s, sz := t.subtreeSum(c)
		sum += s
		size += sz
	}
	return sum, size
}

func (t *Tree) rotate(x int) {
	if len(t.children[x]) == 0 {
		return
	}
	maxSize := -1
	hvIdx := -1
	hv := 0
	for idx, c := range t.children[x] {
		_, sz := t.subtreeSum(c)
		if sz > maxSize || (sz == maxSize && c < hv) {
			maxSize = sz
			hv = c
			hvIdx = idx
		}
	}
	if hvIdx == -1 {
		return
	}
	p := t.parent[x]
	t.children[x] = append(t.children[x][:hvIdx], t.children[x][hvIdx+1:]...)
	if p != 0 {
		for i, c := range t.children[p] {
			if c == x {
				t.children[p][i] = hv
				break
			}
		}
	}
	t.parent[hv] = p
	t.parent[x] = hv
	t.children[hv] = append(t.children[hv], x)
}

func naive(n, m int, vals []int64, edges [][2]int, ops [][2]int) []int64 {
	tr := newTree(n)
	copy(tr.value[1:], vals)
	tr.orient(1, edges)
	res := []int64{}
	for _, op := range ops {
		if op[0] == 1 {
			s, _ := tr.subtreeSum(op[1])
			res = append(res, s)
		} else {
			tr.rotate(op[1])
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierD.go path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if b, err := filepath.Abs(bin); err == nil {
		bin = b
	}

	rand.Seed(4)
	const T = 100
	for tc := 0; tc < T; tc++ {
		n := rand.Intn(4) + 2
		m := rand.Intn(4) + 1
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			vals[i] = int64(rand.Intn(5) + 1)
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u := i + 2
			v := rand.Intn(i+1) + 1
			edges[i] = [2]int{u, v}
		}
		ops := make([][2]int, m)
		for i := 0; i < m; i++ {
			typ := rand.Intn(2) + 1
			x := rand.Intn(n) + 1
			ops[i] = [2]int{typ, x}
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(vals[i], 10))
		}
		input.WriteByte('\n')
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		for _, op := range ops {
			input.WriteString(fmt.Sprintf("%d %d\n", op[0], op[1]))
		}

		expected := naive(n, m, vals, edges, ops)
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Printf("test %d binary error: %v\n", tc+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expected) {
			fmt.Printf("test %d wrong number of outputs\n", tc+1)
			os.Exit(1)
		}
		for i, exp := range expected {
			got, err := strconv.ParseInt(fields[i], 10, 64)
			if err != nil || got != exp {
				fmt.Printf("test %d failed at output %d: expected %d got %s\n", tc+1, i+1, exp, fields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
