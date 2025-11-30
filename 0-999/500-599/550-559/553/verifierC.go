package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
2 0
2 1
1 2 1
6 6
5 1 0
4 6 1
6 5 1
5 4 1
1 3 1
3 4 1
6 5
5 2 0
2 1 0
3 2 0
5 3 0
4 6 1
6 11
3 4 0
4 6 1
6 5 0
4 3 1
5 3 1
4 3 1
4 6 0
3 6 0
5 3 1
3 6 1
3 6 0
5 8
3 5 0
3 1 0
1 5 0
3 5 0
1 5 0
3 2 0
1 4 0
1 3 1
3 1
3 1 0
2 0
2 0
2 1
2 1 0
3 0
5 9
1 2 0
1 3 0
3 4 0
3 4 0
3 4 0
4 2 0
3 1 0
4 2 1
4 5 1
3 2
2 3 1
3 1 0
2 1
1 2 0
6 1
2 6 1
2 1
1 2 1
5 4
5 1 0
1 4 1
2 1 0
2 1 0
2 0
3 0
3 0
6 14
4 3 1
2 6 0
6 4 1
5 1 0
4 5 0
1 6 1
3 1 0
5 3 1
6 3 1
1 6 1
1 3 0
6 1 1
1 4 1
4 2 0
2 1
1 2 1
2 0
5 3
1 5 1
4 2 1
4 1 1
2 0
2 1
2 1 0
2 1
1 2 1
4 3
2 3 1
4 3 1
2 4 1
6 13
6 1 0
1 3 0
5 2 1
1 6 0
2 3 1
2 6 1
4 2 1
1 2 1
1 3 0
6 5 1
2 4 0
4 3 0
4 6 0
6 12
6 2 1
5 1 1
3 4 1
6 4 1
3 5 1
6 1 0
5 2 1
6 4 0
3 4 1
6 2 0
1 4 0
6 5 1
3 0
5 8
2 3 0
4 5 0
1 5 1
4 3 0
4 2 0
3 1 1
5 1 1
3 4 1
6 14
1 5 1
2 4 1
6 2 0
2 4 1
4 6 1
2 1 1
5 4 0
3 1 1
5 4 0
6 3 1
6 2 1
6 5 1
1 3 1
3 6 1
5 8
1 5 0
4 5 0
5 1 1
1 2 1
5 2 1
1 4 1
2 3 1
1 3 1
4 1
4 3 1
3 3
3 1 1
2 1 0
1 2 0
4 1
3 2 0
3 2
3 2 1
3 2 1
6 10
6 4 1
5 6 0
3 4 1
2 3 1
4 1 0
3 4 0
1 3 0
3 1 1
1 5 1
2 5 1
6 9
4 6 0
3 2 1
1 2 0
3 2 1
3 1 1
2 6 0
6 5 0
3 6 0
2 1 1
5 4
2 3 0
3 5 1
2 1 1
4 5 1
6 2
5 6 1
4 6 1
3 3
2 3 1
1 2 0
1 3 0
2 1
2 1 1
2 1
1 2 0
4 6
4 1 1
2 1 1
2 1 0
1 2 0
3 2 0
3 4 0
3 0
3 0
5 3
1 5 0
4 3 1
5 2 0
3 1
1 2 0
4 2
4 2 1
2 4 1
4 6
1 4 1
3 4 1
2 4 0
1 2 0
2 4 0
2 1 1
5 8
1 2 0
5 3 0
2 1 0
5 2 0
4 3 0
2 3 1
5 1 1
4 1 1
4 5
1 3 0
2 4 1
3 4 0
4 2 0
4 3 1
5 9
3 5 0
5 4 1
2 4 1
3 1 0
3 5 1
2 5 0
3 2 1
1 2 1
1 3 1
4 6
4 1 1
4 1 1
2 4 0
4 3 0
4 2 0
3 2 0
6 10
5 6 1
5 4 0
6 2 0
3 6 0
2 5 1
4 1 1
4 3 0
2 1 0
5 6 0
1 4 1
5 3
5 3 0
3 5 1
5 4 0
5 3
3 1 0
4 1 0
2 3 0
6 1
5 2 1
2 0
3 3
3 1 0
2 3 1
2 1 1
4 5
1 3 0
3 4 1
2 1 0
3 4 0
3 1 0
6 13
4 2 0
2 5 0
1 5 1
6 2 1
1 2 0
4 1 1
1 5 0
2 5 0
1 6 1
1 5 1
5 6 0
1 5 1
4 2 0
5 2
4 3 0
5 2 1
5 9
4 2 1
5 2 1
2 1 1
2 4 1
4 2 1
4 3 0
3 1 1
5 1 1
3 1 1
6 8
6 3 0
4 6 1
3 5 1
6 2 0
4 3 1
4 1 0
6 4 0
4 5 1
6 14
1 6 1
4 5 1
5 3 0
2 6 0
4 6 0
4 3 1
5 4 1
6 3 1
6 1 1
5 1 0
4 1 0
1 3 1
5 3 1
2 6 1
2 1
2 1 1
2 0
4 0
4 3
2 1 0
3 4 1
1 4 0
3 2
1 2 1
2 3 1
5 4
1 5 0
5 3 0
3 2 0
5 2 0
6 14
2 4 1
5 2 0
1 4 1
6 3 0
6 2 0
4 5 0
6 2 0
6 2 0
5 1 1
6 2 0
3 4 0
1 3 0
3 1 1
4 1 0
5 2
2 4 0
3 2 1
6 2
5 6 1
3 1 0
2 1
1 2 1
3 0
4 3
1 2 0
3 1 0
2 1 1
2 1
1 2 0
6 4
3 2 1
6 2 1
5 1 0
3 6 0
2 0
3 3
2 3 1
2 3 1
2 1 0
2 0
2 0
4 6
2 3 0
1 2 0
3 1 0
4 2 1
1 3 1
4 3 0
6 3
6 1 0
6 4 0
4 6 0
4 4
1 2 0
1 3 0
2 4 1
4 1 1
3 0
4 0
5 0
6 3
6 3 1
5 2 1
4 2 0
3 1
1 3 0
5 2
1 3 1
4 1 0
5 0
2 0
3 3
1 2 1
1 2 1
2 3 0
6 9
3 5 1
4 3 1
2 6 1
3 5 1
5 4 1
5 1 0
3 4 1
6 1 1
5 6 1
3 1
2 1 0`

type edge struct{ a, b, c int }

type testCaseC struct {
	n     int
	edges []edge
}

func parseTestcases() ([]testCaseC, error) {
	in := bufio.NewReader(strings.NewReader(embeddedTestcases))
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseC, T)
	for i := 0; i < T; i++ {
		var n, m int
		if _, err := fmt.Fscan(in, &n, &m); err != nil {
			return nil, err
		}
		edges := make([]edge, m)
		for j := 0; j < m; j++ {
			fmt.Fscan(in, &edges[j].a, &edges[j].b, &edges[j].c)
		}
		cases[i] = testCaseC{n: n, edges: edges}
	}
	return cases, nil
}

type DSU struct {
	parent []int
	parity []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	d := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{parent: p, parity: d}
}
func (d *DSU) Find(x int) (int, int) {
	if d.parent[x] != x {
		r, p := d.Find(d.parent[x])
		d.parent[x] = r
		d.parity[x] ^= p
	}
	return d.parent[x], d.parity[x]
}
func (d *DSU) Union(x, y, val int) bool {
	rx, px := d.Find(x)
	ry, py := d.Find(y)
	if rx == ry {
		return (px ^ py) == val
	}
	d.parent[rx] = ry
	d.parity[rx] = px ^ py ^ val
	return true
}

const mod int = 1000000007

func modPow(a, b int) int {
	res := 1
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveCase(tc testCaseC) string {
	d := NewDSU(tc.n)
	ok := true
	for _, e := range tc.edges {
		if !d.Union(e.a, e.b, e.c^1) {
			ok = false
		}
	}
	if !ok {
		return "0"
	}
	seen := make(map[int]bool)
	comp := 0
	for i := 1; i <= tc.n; i++ {
		r, _ := d.Find(i)
		if !seen[r] {
			seen[r] = true
			comp++
		}
	}
	ans := modPow(2, comp-1)
	return strconv.Itoa(ans)
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.edges)))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.a, e.b, e.c))
		}
		expected := solveCase(tc)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
