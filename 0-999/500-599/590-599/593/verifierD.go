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

const testcasesDRaw = `100
3 1
1 2 10
1 3 2
1 1 3 8
3 4
1 2 5
2 3 9
1 3 3 16
2 1 2
1 1 1 4
1 1 3 18
2 5
1 2 8
2 1 1
2 1 5
1 2 2 13
2 1 4
2 1 4
4 5
1 2 2
2 3 10
3 4 5
2 3 3
1 1 3 11
2 3 4
2 3 1
2 3 3
6 1
1 2 10
1 3 1
1 4 5
4 5 2
5 6 1
1 2 3 17
4 3
1 2 7
2 3 8
2 4 6
1 3 1 11
2 1 4
1 3 3 20
3 6
1 2 9
2 3 8
1 1 1 6
2 2 3
2 2 3
1 3 2 6
2 2 1
2 2 3
4 5
1 2 2
2 3 1
2 4 10
1 3 3 16
1 2 2 12
2 1 1
1 4 1 11
1 1 4 9
3 6
1 2 8
1 3 7
2 2 1
1 1 2 14
2 2 1
1 3 2 9
1 1 3 17
1 3 3 1
5 3
1 2 9
1 3 2
3 4 4
1 5 9
2 2 2
2 4 4
1 3 2 2
4 6
1 2 6
2 3 4
2 4 1
1 1 2 20
1 2 1 2
2 1 5
1 2 2 18
1 4 1 12
1 1 2 19
2 3
1 2 6
1 2 2 3
1 2 2 5
1 1 1 17
2 3
1 2 10
1 2 2 9
1 2 2 12
1 1 2 6
4 4
1 2 6
2 3 5
3 4 6
2 2 2
1 4 4 19
2 1 3
1 4 4 19
2 1
1 2 5
2 1 4
2 6
1 2 7
2 1 4
2 1 4
1 1 1 11
1 2 1 19
2 1 2
1 2 2 16
4 5
1 2 5
1 3 8
3 4 3
1 3 1 10
2 3 2
2 1 5
1 3 2 5
2 2 5
5 3
1 2 10
1 3 1
3 4 10
4 5 7
1 2 3 16
1 5 2 15
2 2 2
5 4
1 2 7
1 3 9
2 4 7
1 5 5
1 2 2 10
1 3 2 19
2 2 1
1 3 5 10
4 2
1 2 5
2 3 2
2 4 5
1 3 1 18
2 1 1
4 6
1 2 10
2 3 6
3 4 10
2 1 4
1 1 3 18
1 2 1 1
1 3 1 12
2 3 5
2 3 4
6 6
1 2 2
1 3 4
2 4 9
4 5 3
1 6 2
2 3 2
2 5 1
1 3 1 4
2 2 3
1 6 6 6
1 4 4 7
3 1
1 2 7
1 3 9
1 3 2 11
4 5
1 2 8
1 3 10
1 4 10
2 1 3
1 1 1 4
1 1 3 4
1 3 1 4
1 2 3 8
3 2
1 2 1
1 3 8
1 3 1 14
2 1 5
2 3
1 2 5
1 2 1 16
1 2 1 15
2 1 5
5 3
1 2 5
2 3 3
1 4 10
4 5 2
2 3 5
1 3 4 7
1 3 5 11
2 3
1 2 7
2 1 5
2 1 1
1 2 2 11
2 5
1 2 2
1 1 2 2
1 1 1 17
2 1 2
2 1 3
1 1 1 5
6 2
1 2 8
2 3 7
1 4 1
4 5 2
2 6 3
1 6 4 3
1 1 2 14
5 4
1 2 10
1 3 9
1 4 9
4 5 3
2 2 1
2 4 5
1 4 5 14
1 5 2 3
4 5
1 2 6
2 3 4
3 4 7
1 4 3 11
2 1 1
1 2 1 8
1 2 3 7
1 4 2 9
6 3
1 2 2
2 3 8
3 4 2
2 5 2
5 6 5
1 4 4 7
1 2 1 14
1 4 6 13
2 3
1 2 10
2 1 1
1 1 1 11
1 2 2 8
2 3
1 2 10
1 1 1 1
1 2 1 13
1 1 1 6
2 5
1 2 6
1 2 1 15
2 1 3
2 1 4
2 1 2
2 1 5
6 1
1 2 9
2 3 6
2 4 4
1 5 6
3 6 2
2 5 5
3 1
1 2 1
2 3 6
1 1 1 17
3 2
1 2 4
2 3 10
2 1 1
1 3 3 12
4 4
1 2 5
2 3 6
1 4 2
2 3 1
2 1 2
1 4 2 18
1 3 3 9
6 6
1 2 9
1 3 7
3 4 2
4 5 8
1 6 2
2 2 4
2 4 5
2 3 4
1 2 1 19
1 3 5 7
2 1 2
5 6
1 2 3
1 3 4
2 4 9
2 5 5
1 3 3 13
1 2 5 10
2 4 3
1 4 4 18
2 1 2
2 3 4
3 3
1 2 9
2 3 10
1 1 2 7
2 2 4
1 1 3 8
6 5
1 2 8
1 3 4
2 4 3
2 5 8
2 6 1
1 2 4 13
1 1 2 20
1 3 6 2
2 3 4
1 5 6 12
2 6
1 2 6
1 1 1 12
2 1 3
1 1 1 3
2 1 2
1 2 2 17
1 1 1 1
6 4
1 2 7
2 3 3
1 4 10
3 5 8
5 6 4
2 4 1
1 3 5 10
2 4 5
1 4 1 16
4 6
1 2 9
2 3 8
3 4 2
2 1 4
1 4 2 16
2 1 3
2 1 5
1 3 1 6
2 3 5
6 4
1 2 5
1 3 1
3 4 7
3 5 5
4 6 1
2 1 1
1 6 1 2
2 4 5
1 3 1 6
5 5
1 2 2
1 3 5
1 4 8
2 5 4
1 2 5 1
1 4 4 11
1 1 1 20
1 1 2 13
2 3 1
5 3
1 2 9
1 3 2
2 4 6
2 5 9
2 2 5
2 3 5
1 2 3 3
2 3
1 2 7
2 1 1
1 2 1 6
2 1 5
2 5
1 2 10
2 1 2
2 1 1
1 2 1 3
1 2 1 4
1 2 2 19
6 6
1 2 3
1 3 6
2 4 5
3 5 1
1 6 7
1 3 3 1
1 3 2 7
1 3 5 3
1 6 5 8
2 1 3
2 5 4
2 6
1 2 4
1 1 2 17
1 2 1 12
2 1 3
2 1 4
2 1 3
1 2 1 5
4 4
1 2 9
2 3 4
1 4 2
2 1 2
1 3 4 3
2 1 5
2 1 3
3 2
1 2 8
2 3 6
1 2 2 14
2 2 2
5 2
1 2 9
2 3 6
1 4 1
2 5 6
1 2 1 10
1 1 1 15
6 6
1 2 1
2 3 8
1 4 4
1 5 4
2 6 10
1 3 1 10
2 1 3
2 3 2
1 6 6 14
2 2 4
1 5 6 11
2 3
1 2 2
1 1 1 17
1 1 2 4
1 2 2 19
5 1
1 2 10
1 3 3
3 4 3
4 5 3
1 1 5 2
6 4
1 2 7
2 3 8
2 4 10
1 5 3
5 6 1
1 4 3 14
1 1 1 11
1 2 5 19
1 6 2 19
4 2
1 2 5
1 3 7
3 4 2
2 2 4
1 3 4 5
5 5
1 2 2
2 3 7
2 4 6
4 5 6
2 4 3
1 4 3 7
1 1 4 14
1 4 2 14
1 1 2 10
2 2
1 2 2
1 1 2 11
1 1 1 17
5 1
1 2 8
1 3 1
2 4 8
3 5 4
2 2 3
4 2
1 2 6
2 3 5
2 4 8
1 4 2 20
1 3 4 14
5 1
1 2 8
2 3 1
2 4 8
3 5 8
2 1 5
4 6
1 2 10
2 3 6
3 4 8
2 2 4
2 2 5
1 3 1 12
2 1 5
1 4 2 2
2 3 5
3 5
1 2 2
2 3 2
2 2 3
1 2 3 5
2 1 5
2 2 3
1 3 3 18
3 3
1 2 4
1 3 8
2 1 3
2 1 5
2 1 1
5 2
1 2 10
1 3 7
1 4 8
3 5 10
1 1 3 7
1 3 1 15
4 3
1 2 7
2 3 4
3 4 9
1 1 3 13
1 2 4 20
2 2 5
4 3
1 2 3
1 3 3
1 4 5
1 2 3 13
2 3 1
2 3 5
5 1
1 2 9
2 3 1
2 4 7
2 5 2
2 3 4
6 3
1 2 8
2 3 3
3 4 10
3 5 9
5 6 8
2 3 1
2 1 3
2 4 2
6 3
1 2 2
2 3 3
2 4 5
2 5 3
5 6 2
2 5 2
1 1 5 12
2 4 3
3 1
1 2 8
1 3 10
2 2 3
2 3
1 2 2
2 1 1
2 1 1
2 1 3
3 3
1 2 4
2 3 6
2 1 3
2 2 2
1 3 3 19
4 4
1 2 8
1 3 3
2 4 1
2 1 4
2 1 1
1 2 1 4
2 1 5
2 2
1 2 3
1 2 1 15
1 1 1 17
2 1
1 2 10
1 1 2 15
6 5
1 2 5
2 3 10
3 4 10
2 5 6
3 6 3
2 4 2
1 6 2 8
1 5 4 16
1 5 2 20
2 3 1
2 4
1 2 3
1 2 1 7
1 2 2 14
2 1 3
2 1 3
5 4
1 2 10
1 3 1
1 4 10
4 5 8
1 4 5 1
1 1 2 11
2 3 1
2 3 2
3 2
1 2 1
1 3 1
1 2 3 4
2 2 2
3 4
1 2 8
1 3 10
1 3 3 6
2 2 5
1 3 2 14
2 1 2
2 6
1 2 4
1 1 1 5
2 1 3
2 1 1
2 1 2
1 1 2 20
1 2 2 18
3 1
1 2 10
1 3 10
2 2 3
6 5
1 2 4
2 3 4
3 4 8
1 5 7
5 6 8
1 3 2 11
2 1 3
2 4 1
2 5 4
2 3 5
6 2
1 2 1
2 3 2
3 4 4
4 5 5
5 6 5
2 1 5
1 6 1 1
6 1
1 2 1
1 3 2
1 4 5
3 5 2
1 6 2
2 5 5
2 6
1 2 6
2 1 1
2 1 4
1 1 1 12
1 2 1 14
1 1 1 9
2 1 3
3 5
1 2 7
2 3 4
1 3 1 15
2 1 2
2 2 3
2 1 5
1 2 2 2
4 5
1 2 10
2 3 4
3 4 5
1 3 3 13
1 1 4 7
2 1 2
1 3 4 15
2 2 2
6 4
1 2 2
1 3 10
2 4 5
1 5 1
1 6 2
2 2 5
1 3 5 4
1 3 6 16
1 1 3 18
6 4
1 2 10
2 3 3
3 4 9
2 5 1
4 6 10
1 2 5 8
2 3 1
1 4 3 5
2 4 2
5 2
1 2 1
2 3 8
3 4 8
1 5 6
1 2 3 2
2 3 3
3 1
1 2 2
2 3 6
2 2 5
2 1
1 2 4
1 1 2 11
`

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesDRaw)
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
