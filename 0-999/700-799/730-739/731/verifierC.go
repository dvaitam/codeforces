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

type DSU struct {
	parent []int
	rank   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), rank: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(x, y int) {
	rx := d.Find(x)
	ry := d.Find(y)
	if rx == ry {
		return
	}
	if d.rank[rx] < d.rank[ry] {
		d.parent[rx] = ry
	} else if d.rank[ry] < d.rank[rx] {
		d.parent[ry] = rx
	} else {
		d.parent[ry] = rx
		d.rank[rx]++
	}
}

func solveC(n, m, k int, colors []int, pairs [][2]int) string {
	dsu := NewDSU(n)
	used := make([]bool, n)
	for _, p := range pairs {
		l := p[0]
		r := p[1]
		dsu.Union(l, r)
		used[l] = true
		used[r] = true
	}
	compCount := make(map[int]int)
	maxFreq := make(map[int]int)
	colorCount := make(map[int]map[int]int)
	for i := 0; i < n; i++ {
		if !used[i] {
			continue
		}
		root := dsu.Find(i)
		compCount[root]++
		col := colors[i]
		cmap := colorCount[root]
		if cmap == nil {
			cmap = make(map[int]int)
			colorCount[root] = cmap
		}
		cmap[col]++
		if cmap[col] > maxFreq[root] {
			maxFreq[root] = cmap[col]
		}
	}
	result := 0
	for root, cnt := range compCount {
		result += cnt - maxFreq[root]
	}
	return fmt.Sprintf("%d\n", result)
}

func genCaseC(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 2
	k := rng.Intn(5) + 1
	m := rng.Intn(n)
	colors := make([]int, n)
	for i := range colors {
		colors[i] = rng.Intn(k) + 1
	}
	pairs := make([][2]int, m)
	for i := 0; i < m; i++ {
		l := rng.Intn(n)
		r := rng.Intn(n)
		for r == l {
			r = rng.Intn(n)
		}
		pairs[i] = [2]int{l, r}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i, c := range colors {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c)
	}
	sb.WriteByte('\n')
	for _, p := range pairs {
		fmt.Fprintf(&sb, "%d %d\n", p[0]+1, p[1]+1)
	}
	return sb.String(), solveC(n, m, k, colors, pairs)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseC(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
