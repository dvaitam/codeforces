package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type edge struct {
	u, v int
	a, b int64
	idx  int
}

type dsu struct {
	parent []int
	rank   []int
}

func newDSU(n int) *dsu {
	d := &dsu{parent: make([]int, n), rank: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
	}
	return d
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) union(a, b int) bool {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return false
	}
	if d.rank[ra] < d.rank[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	if d.rank[ra] == d.rank[rb] {
		d.rank[ra]++
	}
	return true
}

func solve(n, m int, k int64, edges []edge) int64 {
	used := make([]int64, m)
	var total int64
	for step := int64(0); step < k; step++ {
		weights := make([]struct {
			w int64
			e edge
		}, m)
		for i, e := range edges {
			w := e.a*(2*used[i]+1) + e.b
			weights[i] = struct {
				w int64
				e edge
			}{w, e}
		}
		sort.Slice(weights, func(i, j int) bool { return weights[i].w < weights[j].w })
		d := newDSU(n)
		cnt := 0
		for _, item := range weights {
			if d.union(item.e.u, item.e.v) {
				used[item.e.idx]++
				total += item.w
				cnt++
				if cnt == n-1 {
					break
				}
			}
		}
	}
	return total
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesI.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		kVal, _ := strconv.ParseInt(fields[2], 10, 64)
		edges := make([]edge, m)
		pos := 3
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			a, _ := strconv.ParseInt(fields[pos+2], 10, 64)
			b, _ := strconv.ParseInt(fields[pos+3], 10, 64)
			edges[i] = edge{u - 1, v - 1, a, b, i}
			pos += 4
		}
		exp := solve(n, m, kVal, edges)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, m, kVal))
		for i := 0; i < m; i++ {
			e := edges[i]
			input.WriteString(fmt.Sprintf("%d %d %d %d\n", e.u+1, e.v+1, e.a, e.b))
		}

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != fmt.Sprintf("%d", exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
