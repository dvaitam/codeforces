package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// DSU structure
type DSU struct{ parent, rank []int }

func NewDSU(n int) *DSU {
	p := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	return &DSU{p, r}
}
func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}
func (d *DSU) Union(x, y int) {
	x = d.Find(x)
	y = d.Find(y)
	if x == y {
		return
	}
	if d.rank[x] < d.rank[y] {
		x, y = y, x
	}
	d.parent[y] = x
	if d.rank[x] == d.rank[y] {
		d.rank[x]++
	}
}

func canSort(p []int, l, r int) bool {
	n := len(p)
	d := NewDSU(n + 1)
	for x := 1; x <= n; x++ {
		for y := x + 1; y <= n; y++ {
			s := x + y
			if s >= l && s <= r {
				d.Union(x, y)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if d.Find(i) != d.Find(p[i-1]) {
			return false
		}
	}
	return true
}

func solve(p []int) int {
	n := len(p)
	cnt := 0
	for l := 1; l <= 2*n; l++ {
		for r := l; r <= 2*n; r++ {
			if canSort(p, l, r) {
				cnt++
			}
		}
	}
	return cnt
}

type Case struct {
	p   []int
	ans int
}

func genCases(n int) []Case {
	rand.Seed(time.Now().UnixNano())
	cs := make([]Case, n)
	for i := 0; i < n; i++ {
		size := rand.Intn(4) + 2
		p := rand.Perm(size)
		for j := 0; j < size; j++ {
			p[j]++
		}
		cs[i] = Case{p, solve(p)}
	}
	return cs
}

func buildInput(cs []Case) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintln(&sb, len(c.p))
		for j, v := range c.p {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cs := genCases(100)
	input := buildInput(cs)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cs) {
		fmt.Printf("expected %d outputs got %d\n", len(cs), len(outputs))
		os.Exit(1)
	}
	for i, res := range outputs {
		v, err := strconv.Atoi(res)
		if err != nil || v != cs[i].ans {
			fmt.Printf("mismatch on case %d: expected %d got %s\n", i+1, cs[i].ans, res)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
