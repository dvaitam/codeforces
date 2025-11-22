package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1_000_000_007

type dsu struct {
	parent []int
	sz     []int
	edges  []int
}

func newDSU(n int) *dsu {
	p := make([]int, n)
	s := make([]int, n)
	e := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &dsu{parent: p, sz: s, edges: e}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.sz[ra] += d.sz[rb]
	d.edges[ra] += d.edges[rb]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func solveCase(n, m, k int, odds [][2]int) int64 {
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

	oddSet := make(map[int]struct{}, k+1)
	for _, p := range odds {
		id := (p[0]-1)*m + (p[1] - 1)
		oddSet[id] = struct{}{}
	}

	type variable struct {
		a, b     int
		len      int
		resolved bool
	}

	vars := make([]variable, k)
	cellToVars := make(map[int][]int)
	ok := true

	for i := 0; i < k; i++ {
		a := odds[i]
		b := odds[i+1]
		list := make([]int, 0, 2)
		for _, d := range dirs {
			x := a[0] + d[0]
			y := a[1] + d[1]
			if x < 1 || x > n || y < 1 || y > m {
				continue
			}
			if abs(x-b[0])+abs(y-b[1]) != 1 {
				continue
			}
			id := (x-1)*m + (y - 1)
			if _, exists := oddSet[id]; exists {
				continue
			}
			list = append(list, id)
		}
		if len(list) == 0 {
			ok = false
			break
		}
		if len(list) == 1 {
			vars[i] = variable{a: list[0], len: 1}
		} else {
			vars[i] = variable{a: list[0], b: list[1], len: 2}
		}
		for _, id := range list {
			cellToVars[id] = append(cellToVars[id], i)
		}
	}

	if !ok {
		return 0
	}

	used := make(map[int]bool)
	queue := make([]int, 0)
	for i, v := range vars {
		if v.len == 1 {
			queue = append(queue, i)
		}
	}

	for head := 0; head < len(queue) && ok; head++ {
		idx := queue[head]
		v := &vars[idx]
		if v.resolved || v.len != 1 {
			continue
		}
		cell := v.a
		if used[cell] {
			ok = false
			break
		}
		used[cell] = true
		v.resolved = true

		for _, other := range cellToVars[cell] {
			if vars[other].resolved {
				continue
			}
			if vars[other].len == 1 {
				if vars[other].a == cell {
					ok = false
					break
				}
			} else if vars[other].len == 2 {
				if vars[other].a == cell {
					vars[other].a = vars[other].b
					vars[other].len = 1
					queue = append(queue, other)
				} else if vars[other].b == cell {
					vars[other].len = 1
					queue = append(queue, other)
				}
			}
		}
	}

	if !ok {
		return 0
	}

	for i := 0; i < k; i++ {
		if vars[i].len == 0 {
			return 0
		}
	}

	cellID := make(map[int]int)
	getID := func(id int) int {
		if v, exists := cellID[id]; exists {
			return v
		}
		idx := len(cellID)
		cellID[id] = idx
		return idx
	}

	edges := make([][2]int, 0)
	for i := 0; i < k; i++ {
		if vars[i].resolved {
			continue
		}
		if vars[i].len == 2 {
			edges = append(edges, [2]int{vars[i].a, vars[i].b})
			getID(vars[i].a)
			getID(vars[i].b)
		} else if vars[i].len == 1 {
			if used[vars[i].a] {
				return 0
			}
			used[vars[i].a] = true
			vars[i].resolved = true
		} else {
			return 0
		}
	}

	if len(edges) == 0 {
		return 1
	}

	d := newDSU(len(cellID))
	for _, e := range edges {
		u := getID(e[0])
		v := getID(e[1])
		d.union(u, v)
	}
	for _, e := range edges {
		u := getID(e[0])
		root := d.find(u)
		d.edges[root]++
	}

	ans := int64(1)
	for i := 0; i < len(d.parent); i++ {
		if d.find(i) != i {
			continue
		}
		e := d.edges[i]
		vCount := d.sz[i]
		if e > vCount {
			return 0
		}
		if e == vCount {
			ans = ans * 2 % mod
		} else if e == vCount-1 {
			ans = ans * int64(vCount) % mod
		} else {
			return 0
		}
	}

	return ans % mod
}

type test struct {
	raw   string
	cases [][2]int // stores (n*m, placeholder) unused; mainly to keep count
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(2097))
	var tests []test

	// Simple minimal cases.
	tests = append(tests, test{raw: "1\n1 3 1\n1 1\n1 3\n"})
	tests = append(tests, test{raw: "1\n2 2 1\n1 1\n2 2\n"})

	for len(tests) < 120 {
		t := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", t)
		for i := 0; i < t; i++ {
			n := rng.Intn(6) + 1
			m := rng.Intn(6) + 1
			if n*m < 3 {
				i--
				continue
			}
			maxK := (n*m - 1) / 2
			if maxK < 1 {
				i--
				continue
			}
			k := rng.Intn(maxK) + 1
			used := make(map[int]struct{})
			odds := make([][2]int, k+1)
			for j := 0; j <= k; j++ {
				for {
					x := rng.Intn(n) + 1
					y := rng.Intn(m) + 1
					id := (x-1)*m + (y - 1)
					if _, ok := used[id]; ok {
						continue
					}
					used[id] = struct{}{}
					odds[j] = [2]int{x, y}
					break
				}
			}
			fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
			for _, p := range odds {
				fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
			}
		}
		tests = append(tests, test{raw: sb.String()})
	}

	return tests
}

func runBinary(bin, input string) (string, error) {
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

func parseAndJudge(input, output string) error {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return fmt.Errorf("failed to read t from input")
	}
	exp := make([]int64, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		var n, m, k int
		fmt.Fscan(in, &n, &m, &k)
		odds := make([][2]int, k+1)
		for i := 0; i <= k; i++ {
			fmt.Fscan(in, &odds[i][0], &odds[i][1])
		}
		exp[caseIdx] = solveCase(n, m, k, odds)
	}

	out := bufio.NewReader(strings.NewReader(output))
	for i := 0; i < t; i++ {
		var got int64
		if _, err := fmt.Fscan(out, &got); err != nil {
			return fmt.Errorf("missing output for test %d", i+1)
		}
		got %= mod
		if got < 0 {
			got += mod
		}
		if got != exp[i] {
			return fmt.Errorf("test %d: expected %d, got %d", i+1, exp[i], got)
		}
	}
	var extra int
	if _, err := fmt.Fscan(out, &extra); err == nil {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.raw)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := parseAndJudge(t.raw, got); err != nil {
			fmt.Printf("Wrong answer on test %d: %v\nInput:\n%s\nOutput:\n%s\n", i+1, err, t.raw, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
