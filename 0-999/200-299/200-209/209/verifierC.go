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
	for i := 0; i <= n; i++ {
		p[i] = i
	}
	return &DSU{p: p}
}

func (d *DSU) Find(x int) int {
	if d.p[x] != x {
		d.p[x] = d.Find(d.p[x])
	}
	return d.p[x]
}

func (d *DSU) Union(x, y int) {
	rx, ry := d.Find(x), d.Find(y)
	if rx != ry {
		d.p[ry] = rx
	}
}

func solve(n int, edges [][2]int) int {
	dsu := NewDSU(n)
	deg := make([]int, n+1)
	for _, e := range edges {
		x, y := e[0], e[1]
		deg[x]++
		deg[y]++
		dsu.Union(x, y)
	}
	hasEdge := make([]bool, n+1)
	for v := 1; v <= n; v++ {
		if deg[v] > 0 {
			r := dsu.Find(v)
			hasEdge[r] = true
		}
	}
	oddMap := make([]int, n+1)
	totalOdd := 0
	for v := 1; v <= n; v++ {
		if deg[v]%2 == 1 {
			totalOdd++
			r := dsu.Find(v)
			oddMap[r]++
		}
	}
	root1 := dsu.Find(1)
	nonZero := []int{}
	zeroCnt := 0
	for r := 1; r <= n; r++ {
		if r == root1 || !hasEdge[r] {
			continue
		}
		if oddMap[r] > 0 {
			nonZero = append(nonZero, oddMap[r])
		} else {
			zeroCnt++
		}
	}
	odd1 := oddMap[root1]
	connections := 0
	for range nonZero {
		if odd1 > 0 {
			odd1--
			totalOdd -= 2
		} else {
			odd1++
		}
		connections++
	}
	for i := 0; i < zeroCnt; i++ {
		if odd1 > 0 {
		} else {
			odd1++
			totalOdd += 2
		}
		connections++
	}
	needed := totalOdd / 2
	return connections + needed
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			x, _ := strconv.Atoi(parts[2+2*i])
			y, _ := strconv.Atoi(parts[2+2*i+1])
			edges[i] = [2]int{x, y}
		}
		exp := solve(n, edges)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d\n", n, m)
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.Atoi(gotStr)
		if err2 != nil || got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
