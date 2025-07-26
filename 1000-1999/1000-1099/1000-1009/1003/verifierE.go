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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func buildTree(N, d, k int) ([][2]int, bool) {
	if d >= N || (k == 1 && N > 2) {
		return nil, false
	}
	if N == 1 && d == 1 {
		return [][2]int{}, true
	}
	count := make([]int, N+2)
	depth := make([]int, N+2)
	edges := make([][2]int, 0, N-1)
	for i := 1; i <= d; i++ {
		u := i
		v := i + 1
		edges = append(edges, [2]int{u, v})
		count[u]++
		count[v]++
		depth[u] = min(u-1, d-u+1)
		depth[v] = min(u, d-u)
	}
	i := d + 2
	j := 2
	for i <= N {
		for j < i && (count[j] == k || depth[j] == 0) {
			j++
		}
		if i == j {
			return nil, false
		}
		edges = append(edges, [2]int{i, j})
		count[i]++
		count[j]++
		depth[i] = depth[j] - 1
		i++
	}
	if len(edges) != N-1 {
		return nil, false
	}
	return edges, true
}

func diameter(n int, edges [][2]int) int {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	bfs := func(start int) (int, int) {
		dist := make([]int, n+1)
		for i := range dist {
			dist[i] = -1
		}
		q := []int{start}
		dist[start] = 0
		idx := 0
		for idx < len(q) {
			v := q[idx]
			idx++
			for _, to := range g[v] {
				if dist[to] == -1 {
					dist[to] = dist[v] + 1
					q = append(q, to)
				}
			}
		}
		far := start
		for i, dv := range dist {
			if dv > dist[far] {
				far = i
			}
		}
		return far, dist[far]
	}
	v, _ := bfs(1)
	w, dist := bfs(v)
	_ = w
	return dist
}

func checkTree(N, d, k int, edges [][2]int) error {
	if len(edges) != N-1 {
		return fmt.Errorf("expected %d edges got %d", N-1, len(edges))
	}
	deg := make([]int, N+1)
	parent := make([]int, N+1)
	count := 0
	for _, e := range edges {
		u, v := e[0], e[1]
		if u < 1 || u > N || v < 1 || v > N {
			return fmt.Errorf("invalid vertex")
		}
		deg[u]++
		deg[v]++
		if deg[u] > k || deg[v] > k {
			return fmt.Errorf("degree limit exceeded")
		}
		// union find for connectivity
		pu, pv := u, v
		for parent[pu] != 0 {
			pu = parent[pu]
		}
		for parent[pv] != 0 {
			pv = parent[pv]
		}
		if pu != pv {
			parent[pv] = pu
			count++
		}
	}
	if count != N-1 {
		return fmt.Errorf("graph not connected")
	}
	if diameter(N, edges) != d {
		return fmt.Errorf("diameter mismatch")
	}
	return nil
}

func runCase(bin string, n, d, k int) error {
	input := fmt.Sprintf("%d %d %d\n", n, d, k)
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	first := scanner.Text()
	edgesExp, _ := buildTree(n, d, k)
	if first == "NO" {
		if len(edgesExp) != 0 {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("invalid output")
	}
	if len(edgesExp) == 0 {
		return fmt.Errorf("expected NO but got YES")
	}
	edges := make([][2]int, 0, n-1)
	for scanner.Scan() {
		uStr := scanner.Text()
		if !scanner.Scan() {
			return fmt.Errorf("incomplete edge")
		}
		vStr := scanner.Text()
		u, _ := strconv.Atoi(uStr)
		v, _ := strconv.Atoi(vStr)
		edges = append(edges, [2]int{u, v})
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if err := checkTree(n, d, k, edges); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("could not open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		d, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		if err := runCase(bin, n, d, k); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
