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

type Node struct {
	comp   int
	parent [6]int16
	active [6]bool
}

func find(parent []int16, x int) int {
	for int(parent[x]) != x {
		x = int(parent[x])
	}
	return x
}

func union(parent []int16, x, y int) bool {
	fx := find(parent, x)
	fy := find(parent, y)
	if fx == fy {
		return false
	}
	parent[fy] = int16(fx)
	return true
}

func buildLeaf(cols []string, idx int) Node {
	var n Node
	for i := 0; i < 6; i++ {
		n.parent[i] = int16(i)
	}
	for r := 0; r < 3; r++ {
		if cols[r][idx] == '1' {
			n.active[r] = true
			n.active[r+3] = true
			union(n.parent[:], r, r+3)
			if r > 0 && cols[r-1][idx] == '1' {
				union(n.parent[:], r-1, r)
				union(n.parent[:], r-1+3, r+3)
			}
		}
	}
	vis := make(map[int]bool)
	for r := 0; r < 3; r++ {
		if n.active[r] {
			rt := find(n.parent[:], r)
			if !vis[rt] {
				vis[rt] = true
				n.comp++
			}
		}
	}
	return n
}

func merge(a, b Node) Node {
	parent := make([]int16, 12)
	for i := range parent {
		parent[i] = int16(i)
	}
	active := make([]bool, 12)
	for i := 0; i < 6; i++ {
		active[i] = a.active[i]
		active[6+i] = b.active[i]
	}
	for i := 0; i < 6; i++ {
		if !active[i] {
			continue
		}
		for j := i + 1; j < 6; j++ {
			if !active[j] {
				continue
			}
			if find(a.parent[:], i) == find(a.parent[:], j) {
				union(parent, i, j)
			}
		}
	}
	for i := 0; i < 6; i++ {
		if !active[6+i] {
			continue
		}
		for j := i + 1; j < 6; j++ {
			if !active[6+j] {
				continue
			}
			if find(b.parent[:], i) == find(b.parent[:], j) {
				union(parent, 6+i, 6+j)
			}
		}
	}
	comp := a.comp + b.comp
	for r := 0; r < 3; r++ {
		if a.active[3+r] && b.active[r] {
			if union(parent, 3+r, 6+r) {
				comp--
			}
		}
	}
	var res Node
	res.comp = comp
	for i := 0; i < 6; i++ {
		res.parent[i] = int16(i)
	}
	for i := 0; i < 3; i++ {
		res.active[i] = a.active[i]
	}
	for i := 0; i < 3; i++ {
		res.active[3+i] = b.active[3+i]
	}
	idxMap := []int{0, 1, 2, 9, 10, 11}
	for i := 0; i < 6; i++ {
		if !res.active[i] {
			continue
		}
		for j := i + 1; j < 6; j++ {
			if !res.active[j] {
				continue
			}
			if find(parent, idxMap[i]) == find(parent, idxMap[j]) {
				union(res.parent[:], i, j)
			}
		}
	}
	return res
}

func build(seg []Node, cols []string, id, l, r int) {
	if l == r {
		seg[id] = buildLeaf(cols, l)
		return
	}
	mid := (l + r) / 2
	build(seg, cols, id*2, l, mid)
	build(seg, cols, id*2+1, mid+1, r)
	seg[id] = merge(seg[id*2], seg[id*2+1])
}

func query(seg []Node, id, l, r, L, R int) Node {
	if L <= l && r <= R {
		return seg[id]
	}
	mid := (l + r) / 2
	if R <= mid {
		return query(seg, id*2, l, mid, L, R)
	}
	if L > mid {
		return query(seg, id*2+1, mid+1, r, L, R)
	}
	left := query(seg, id*2, l, mid, L, mid)
	right := query(seg, id*2+1, mid+1, r, mid+1, R)
	return merge(left, right)
}

func solveCase(n int, rows []string, queries [][2]int) []int {
	seg := make([]Node, 4*n)
	build(seg, rows, 1, 0, n-1)
	ans := make([]int, len(queries))
	for i, q := range queries {
		l := q[0] - 1
		r := q[1] - 1
		res := query(seg, 1, 0, n-1, l, r)
		ans[i] = res.comp
	}
	return ans
}

type testCase struct {
	input    string
	expected string
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCase, 0, 100)
	// simple case
	cases = append(cases, func() testCase {
		n := 1
		rows := []string{"1", "1", "1"}
		queries := [][2]int{{1, 1}}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < 3; i++ {
			sb.WriteString(rows[i])
			sb.WriteByte('\n')
		}
		sb.WriteString("1\n1 1\n")
		res := solveCase(n, rows, queries)
		exp := fmt.Sprintf("%d", res[0])
		return testCase{input: sb.String(), expected: exp}
	}())
	for len(cases) < 100 {
		n := rng.Intn(10) + 1
		rows := make([]string, 3)
		for i := 0; i < 3; i++ {
			var sb strings.Builder
			for j := 0; j < n; j++ {
				if rng.Intn(2) == 1 {
					sb.WriteByte('1')
				} else {
					sb.WriteByte('0')
				}
			}
			rows[i] = sb.String()
		}
		q := rng.Intn(10) + 1
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = [2]int{l, r}
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < 3; i++ {
			sb.WriteString(rows[i])
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", queries[i][0], queries[i][1]))
		}
		res := solveCase(n, rows, queries)
		var expBuf strings.Builder
		for i, v := range res {
			if i > 0 {
				expBuf.WriteByte('\n')
			}
			expBuf.WriteString(fmt.Sprintf("%d", v))
		}
		cases = append(cases, testCase{input: sb.String(), expected: expBuf.String()})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genTests()
	for i, tc := range cases {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(tc.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, tc.expected, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
