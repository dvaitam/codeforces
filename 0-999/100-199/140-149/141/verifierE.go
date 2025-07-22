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

type DSU struct{ p []int }

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p}
}
func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}
func (d *DSU) Union(x, y int) {
	fx := d.Find(x)
	fy := d.Find(y)
	if fx != fy {
		d.p[fx] = fy
	}
}

type Edge struct {
	u, v int
	s    bool
}

func expected(n, m int, edges []Edge) string {
	if n%2 == 0 {
		return "-1"
	}
	initialUsd := make([]bool, m)
	usd := make([]bool, m)
	dsu := NewDSU(n)
	num := 0
	for i := 0; i < m; i++ {
		e := edges[i]
		if dsu.Find(e.u) != dsu.Find(e.v) {
			dsu.Union(e.u, e.v)
			initialUsd[i] = true
			if e.s {
				num++
			}
		}
	}
	copy(usd, initialUsd)
	if num*2+1 != n {
		dsu = NewDSU(n)
		if num*2+1 < n {
			newUsd := make([]bool, m)
			newNum := 0
			for i := 0; i < m; i++ {
				if edges[i].s && initialUsd[i] {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
					newNum++
				}
			}
			for i := 0; i < m && newNum*2+1 != n; i++ {
				if edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
					newNum++
				}
			}
			for i := 0; i < m; i++ {
				if !edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
				}
			}
			usd = newUsd
			num = newNum
		} else {
			newUsd := make([]bool, m)
			newNum := num
			for i := 0; i < m; i++ {
				if !edges[i].s && initialUsd[i] {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
				}
			}
			for i := 0; i < m && newNum*2+1 != n; i++ {
				if !edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
					newNum--
				}
			}
			for i := 0; i < m; i++ {
				if edges[i].s && dsu.Find(edges[i].u) != dsu.Find(edges[i].v) {
					dsu.Union(edges[i].u, edges[i].v)
					newUsd[i] = true
				}
			}
			usd = newUsd
			num = newNum
		}
	}
	if num*2+1 == n {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n-1))
		cnt := 0
		for i := 0; i < m && cnt < n-1; i++ {
			if usd[i] {
				sb.WriteString(fmt.Sprintf("%d ", i+1))
				cnt++
			}
		}
		return strings.TrimSpace(sb.String())
	}
	return "-1"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("failed to open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	tests := 0
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+3*m {
			fmt.Printf("test %d: invalid line\n", tests+1)
			os.Exit(1)
		}
		edges := make([]Edge, m)
		idx := 2
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(parts[idx])
			idx++
			v, _ := strconv.Atoi(parts[idx])
			idx++
			typ := parts[idx]
			idx++
			edges[i] = Edge{u, v, typ == "S"}
		}
		expect := expected(n, m, edges)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < m; i++ {
			t := 'E'
			if edges[i].s {
				t = 'S'
			}
			input.WriteString(fmt.Sprintf("%d %d %c\n", edges[i].u, edges[i].v, t))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", tests+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", tests+1, expect, got)
			os.Exit(1)
		}
		tests++
	}
	if err := sc.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", tests)
}
