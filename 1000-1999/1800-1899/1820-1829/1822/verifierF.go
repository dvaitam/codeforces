package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runExe(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func bfs(start int, adj [][]int) ([]int, int) {
	n := len(adj) - 1
	dist := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	q := make([]int, 0, n)
	dist[start] = 0
	q = append(q, start)
	far := start
	for head := 0; head < len(q); head++ {
		v := q[head]
		if dist[v] > dist[far] {
			far = v
		}
		for _, to := range adj[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	return dist, far
}

func solveCase(n int, k, c int64, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	dist1, far1 := bfs(1, adj)
	distA, far2 := bfs(far1, adj)
	distB, _ := bfs(far2, adj)
	best := int64(-1 << 63)
	for v := 1; v <= n; v++ {
		ecc := distA[v]
		if distB[v] > ecc {
			ecc = distB[v]
		}
		profit := int64(ecc)*k - int64(dist1[v])*c
		if profit > best {
			best = profit
		}
	}
	return fmt.Sprintf("%d", best)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var tcase int
	fmt.Sscan(scan.Text(), &tcase)
	for idx := 0; idx < tcase; idx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		var k, c int64
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &k)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &c)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &edges[i][0])
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &edges[i][1])
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d %d\n", n, k, c)
		for _, e := range edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		input := sb.String()
		exp := solveCase(n, k, c, edges)
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
