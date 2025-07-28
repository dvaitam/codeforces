package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseG struct {
	n, m  int
	edges [][3]int
	b, e  int
	exp   string
}

type deque struct {
	data       []int
	head, tail int
}

func newDeque(n int) *deque {
	size := 2*n + 10
	d := &deque{data: make([]int, size)}
	d.head = size / 2
	d.tail = d.head
	return d
}

func (d *deque) empty() bool     { return d.head == d.tail }
func (d *deque) pushFront(x int) { d.head--; d.data[d.head] = x }
func (d *deque) pushBack(x int)  { d.data[d.tail] = x; d.tail++ }
func (d *deque) popFront() int   { x := d.data[d.head]; d.head++; return x }

func solveG(n, m int, edges [][3]int, b, e int) string {
	colorID := make(map[int]int)
	vertexColors := make([][]int, n+1)
	colorVerts := [][]int{{}}
	for _, ed := range edges {
		u, v, c := ed[0], ed[1], ed[2]
		id, ok := colorID[c]
		if !ok {
			id = len(colorVerts)
			colorID[c] = id
			colorVerts = append(colorVerts, []int{})
		}
		vertexColors[u] = append(vertexColors[u], id)
		vertexColors[v] = append(vertexColors[v], id)
		colorVerts[id] = append(colorVerts[id], u)
		colorVerts[id] = append(colorVerts[id], v)
	}
	if b == e {
		return "0"
	}
	numColors := len(colorVerts) - 1
	total := n + numColors
	const INF = int(1e9)
	dist := make([]int, total+1)
	for i := 1; i <= total; i++ {
		dist[i] = INF
	}
	dq := newDeque(total)
	dist[b] = 0
	dq.pushFront(b)
	for !dq.empty() {
		v := dq.popFront()
		d := dist[v]
		if v == e {
			break
		}
		if v <= n {
			for _, cid := range vertexColors[v] {
				node := n + cid
				if dist[node] > d+1 {
					dist[node] = d + 1
					dq.pushBack(node)
				}
			}
		} else {
			cid := v - n
			if colorVerts[cid] != nil {
				for _, u := range colorVerts[cid] {
					if dist[u] > d {
						dist[u] = d
						dq.pushFront(u)
					}
				}
				colorVerts[cid] = nil
			}
		}
	}
	return fmt.Sprint(dist[e])
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
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateTests() []testCaseG {
	rng := rand.New(rand.NewSource(7))
	cases := make([]testCaseG, 100)
	for i := range cases {
		n := rng.Intn(5) + 2
		m := rng.Intn(6) + 1
		edges := make([][3]int, m)
		for j := 0; j < m; j++ {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			for v == u {
				v = rng.Intn(n) + 1
			}
			c := rng.Intn(5) + 1
			edges[j] = [3]int{u, v, c}
		}
		b := rng.Intn(n) + 1
		e := rng.Intn(n) + 1
		cases[i] = testCaseG{n: n, m: m, edges: edges, b: b, e: e, exp: solveG(n, m, edges, b, e)}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintln(&sb, 1)
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, ed := range tc.edges {
			fmt.Fprintf(&sb, "%d %d %d\n", ed[0], ed[1], ed[2])
		}
		fmt.Fprintf(&sb, "%d %d\n", tc.b, tc.e)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
