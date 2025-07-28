package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func getPath(u, v int, g [][]int, t int) []int {
	parent := make([]int, t+1)
	for i := range parent {
		parent[i] = -1
	}
	q := []int{u}
	parent[u] = 0
	for len(q) > 0 {
		x := q[0]
		q = q[1:]
		if x == v {
			break
		}
		for _, y := range g[x] {
			if y < 1 || y > t {
				continue
			}
			if parent[y] == -1 {
				parent[y] = x
				q = append(q, y)
			}
		}
	}
	if parent[v] == -1 {
		return nil
	}
	var path []int
	cur := v
	for cur != 0 {
		path = append([]int{cur}, path...)
		cur = parent[cur]
	}
	return path
}

func countCute(t int, g [][]int) int {
	res := 0
	for u := 1; u <= t; u++ {
		for v := u + 1; v <= t; v++ {
			path := getPath(u, v, g, t)
			if len(path) == 0 {
				continue
			}
			minIdx, maxIdx := path[0], path[0]
			for _, x := range path[1:] {
				if x < minIdx {
					minIdx = x
				}
				if x > maxIdx {
					maxIdx = x
				}
			}
			cond1 := u == minIdx
			cond2 := v == maxIdx
			if (cond1 && !cond2) || (!cond1 && cond2) {
				res++
			}
		}
	}
	return res
}

func naive(n int, edges [][2]int, ops []int) []int {
	maxN := 2*n + len(ops) + 5
	g := make([][]int, maxN)
	for _, e := range edges {
		g[e[0]] = append(g[e[0]], e[1])
		g[e[1]] = append(g[e[1]], e[0])
	}
	results := make([]int, 0, len(ops)+1)
	cur := n
	results = append(results, countCute(cur, g))
	for _, p := range ops {
		cur++
		g[p] = append(g[p], cur)
		g[cur] = append(g[cur], p)
		results = append(results, countCute(cur, g))
	}
	return results
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierF.go path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if b, err := filepath.Abs(bin); err == nil {
		bin = b
	}

	rand.Seed(6)
	const T = 100
	for tc := 0; tc < T; tc++ {
		n := rand.Intn(4) + 2
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u := i + 2
			v := rand.Intn(i+1) + 1
			edges[i] = [2]int{u, v}
		}
		m := rand.Intn(3) + 1
		ops := make([]int, m)
		for i := 0; i < m; i++ {
			ops[i] = rand.Intn(n+i) + 1
		}

		var input strings.Builder
		input.WriteString(strconv.Itoa(n) + "\n")
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input.WriteString(strconv.Itoa(m) + "\n")
		for _, k := range ops {
			input.WriteString(strconv.Itoa(k) + "\n")
		}

		expected := naive(n, edges, ops)
		out, err := runBinary(bin, input.String())
		if err != nil {
			fmt.Printf("test %d binary error: %v\n", tc+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expected) {
			fmt.Printf("test %d wrong number of outputs\n", tc+1)
			os.Exit(1)
		}
		for i, exp := range expected {
			got, err := strconv.Atoi(fields[i])
			if err != nil || got != exp {
				fmt.Printf("test %d failed at output %d: expected %d got %s\n", tc+1, i+1, exp, fields[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
