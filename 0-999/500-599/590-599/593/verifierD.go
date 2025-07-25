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

type edge struct {
	to  int
	idx int
}

type inputD struct {
	n, m  int
	edges [][3]int
	ops   [][]int
}

func solve(inp inputD) []int64 {
	n := inp.n
	weights := make([]int64, n)
	g := make([][]edge, n+1)
	for i := 1; i < n; i++ {
		u := inp.edges[i-1][0]
		v := inp.edges[i-1][1]
		w := int64(inp.edges[i-1][2])
		weights[i] = w
		g[u] = append(g[u], edge{v, i})
		g[v] = append(g[v], edge{u, i})
	}
	var res []int64
	for _, op := range inp.ops {
		if op[0] == 1 {
			a, b, y := op[1], op[2], int64(op[3])
			// BFS to find path
			parent := make([]int, n+1)
			pedge := make([]int, n+1)
			for i := 1; i <= n; i++ {
				parent[i] = -1
			}
			q := []int{a}
			parent[a] = 0
			for len(q) > 0 && parent[b] == -1 {
				u := q[0]
				q = q[1:]
				for _, e := range g[u] {
					if parent[e.to] == -1 {
						parent[e.to] = u
						pedge[e.to] = e.idx
						q = append(q, e.to)
					}
				}
			}
			cur := b
			for cur != a {
				y /= weights[pedge[cur]]
				cur = parent[cur]
			}
			res = append(res, y)
		} else {
			p := op[1]
			c := int64(op[2])
			weights[p] = c
		}
	}
	return res
}

func parseInput(data []byte) ([]inputD, error) {
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	tests := make([]inputD, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		edges := make([][3]int, n-1)
		for j := 1; j < n; j++ {
			scan.Scan()
			u, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			w, _ := strconv.Atoi(scan.Text())
			edges[j-1] = [3]int{u, v, w}
		}
		ops := make([][]int, m)
		for j := 0; j < m; j++ {
			scan.Scan()
			tp, _ := strconv.Atoi(scan.Text())
			if tp == 1 {
				scan.Scan()
				a, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				b, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				y, _ := strconv.Atoi(scan.Text())
				ops[j] = []int{1, a, b, y}
			} else {
				scan.Scan()
				p, _ := strconv.Atoi(scan.Text())
				scan.Scan()
				c, _ := strconv.Atoi(scan.Text())
				ops[j] = []int{2, p, c}
			}
		}
		tests[i] = inputD{n, m, edges, ops}
	}
	return tests, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	tests, err := parseInput(data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for ti, tst := range tests {
		expected := solve(tst)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tst.n, tst.m))
		for _, e := range tst.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		for _, op := range tst.ops {
			if op[0] == 1 {
				sb.WriteString(fmt.Sprintf("1 %d %d %d\n", op[1], op[2], op[3]))
			} else {
				sb.WriteString(fmt.Sprintf("2 %d %d\n", op[1], op[2]))
			}
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", ti+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		for i, exp := range expected {
			if !outScan.Scan() {
				fmt.Printf("missing output after query %d in test %d\n", i+1, ti+1)
				os.Exit(1)
			}
			got, _ := strconv.ParseInt(outScan.Text(), 10, 64)
			if got != exp {
				fmt.Printf("test %d query %d expected %d got %d\n", ti+1, i+1, exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
