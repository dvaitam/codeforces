package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type edge struct {
	to int
	w  int64
}

func solveCase(n int, d [][]int64) string {
	const INF = int64(1e18)
	// Basic validation
	for i := 0; i < n; i++ {
		if d[i][i] != 0 {
			return "NO"
		}
		for j := 0; j < n; j++ {
			if d[i][j] != d[j][i] {
				return "NO"
			}
			if i != j && d[i][j] == 0 {
				return "NO"
			}
		}
	}
	if n == 1 {
		return "YES"
	}
	used := make([]bool, n)
	key := make([]int64, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		key[i] = INF
		parent[i] = -1
	}
	key[0] = 0
	for it := 0; it < n; it++ {
		u := -1
		best := INF
		for i := 0; i < n; i++ {
			if !used[i] && key[i] < best {
				best = key[i]
				u = i
			}
		}
		if u == -1 {
			return "NO"
		}
		used[u] = true
		for v := 0; v < n; v++ {
			if !used[v] && d[u][v] < key[v] {
				key[v] = d[u][v]
				parent[v] = u
			}
		}
	}
	adj := make([][]edge, n)
	for v := 1; v < n; v++ {
		u := parent[v]
		if u < 0 || key[v] <= 0 {
			return "NO"
		}
		w := key[v]
		adj[u] = append(adj[u], edge{v, w})
		adj[v] = append(adj[v], edge{u, w})
	}
	dist := make([]int64, n)
	q := make([]int, n)
	for src := 0; src < n; src++ {
		for i := 0; i < n; i++ {
			dist[i] = -1
		}
		head, tail := 0, 0
		q[tail] = src
		tail++
		dist[src] = 0
		for head < tail {
			u := q[head]
			head++
			for _, e := range adj[u] {
				v := e.to
				if dist[v] < 0 {
					dist[v] = dist[u] + e.w
					q[tail] = v
					tail++
				}
			}
		}
		for j := 0; j < n; j++ {
			if dist[j] != d[src][j] {
				return "NO"
			}
		}
	}
	return "YES"
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
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := make([]string, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		d := make([][]int64, n)
		for i := 0; i < n; i++ {
			row := make([]int64, n)
			for j := 0; j < n; j++ {
				scan.Scan()
				val, _ := strconv.ParseInt(scan.Text(), 10, 64)
				row[j] = val
			}
			d[i] = row
		}
		expected[caseIdx] = solveCase(n, d)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
