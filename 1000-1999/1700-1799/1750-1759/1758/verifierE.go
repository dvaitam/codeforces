package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MOD int64 = 1_000_000_007

type DSU struct {
	parent []int
	diff   []int64
	mod    int64
}

func NewDSU(n int, mod int64) *DSU {
	parent := make([]int, n)
	diff := make([]int64, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	return &DSU{parent: parent, diff: diff, mod: mod}
}

func (d *DSU) find(x int) (int, int64) {
	if d.parent[x] == x {
		return x, 0
	}
	r, w := d.find(d.parent[x])
	d.diff[x] = (d.diff[x] + w) % d.mod
	d.parent[x] = r
	return r, d.diff[x]
}

func (d *DSU) unite(x, y int, val int64) bool {
	rx, dx := d.find(x)
	ry, dy := d.find(y)
	if rx == ry {
		if (dx-dy-val)%d.mod != 0 {
			return false
		}
		return true
	}
	d.parent[rx] = ry
	d.diff[rx] = ((val-dx+dy)%d.mod + d.mod) % d.mod
	return true
}

func powMod(a, b, mod int64) int64 {
	res := int64(1)
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func generateTests() []struct {
	n, m int
	h    int64
	grid [][]int64
} {
	r := rand.New(rand.NewSource(5))
	tests := make([]struct {
		n, m int
		h    int64
		grid [][]int64
	}, 100)
	for idx := range tests {
		n := r.Intn(3) + 1
		m := r.Intn(3) + 1
		h := int64(r.Intn(5) + 2)
		grid := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, m)
			for j := 0; j < m; j++ {
				if r.Intn(3) == 0 {
					row[j] = -1
				} else {
					row[j] = int64(r.Intn(int(h)))
				}
			}
			grid[i] = row
		}
		tests[idx] = struct {
			n, m int
			h    int64
			grid [][]int64
		}{n, m, h, grid}
	}
	return tests
}

func solve(t struct {
	n, m int
	h    int64
	grid [][]int64
}) int64 {
	n, m, h := t.n, t.m, t.h
	dsu := NewDSU(n+m, h)
	ok := true
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := t.grid[i][j]
			if x != -1 {
				if !dsu.unite(i, n+j, x%h) {
					ok = false
				}
			}
		}
	}
	if !ok {
		return 0
	}
	comp := 0
	for i := 0; i < n+m; i++ {
		r, _ := dsu.find(i)
		if r == i {
			comp++
		}
	}
	return powMod(h, int64(comp-1), MOD)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tcase := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tcase.n, tcase.m, tcase.h)
		for i := 0; i < tcase.n; i++ {
			for j := 0; j < tcase.m; j++ {
				if j > 0 {
					input.WriteByte(' ')
				}
				fmt.Fprint(&input, tcase.grid[i][j])
			}
			input.WriteByte('\n')
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = &input
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "execution failed: %v\n", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(&out)
	for i, tc := range tests {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "missing output for test %d\n", i+1)
			os.Exit(1)
		}
		gotLine := strings.TrimSpace(scanner.Text())
		got, err := strconv.ParseInt(gotLine, 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid integer on test %d\n", i+1)
			os.Exit(1)
		}
		want := solve(tc)
		if got != want {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output detected")
		os.Exit(1)
	}
	fmt.Println("All test cases passed.")
}
