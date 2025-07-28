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

func expected(n, a, b int, edges [][2]int) string {
	adj := make([][]int, n+1)
	deg := make([]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		deg[u]++
		deg[v]++
	}
	queue := make([]int, 0, n)
	removed := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if deg[i] == 1 {
			queue = append(queue, i)
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		removed[v] = true
		for _, to := range adj[v] {
			if removed[to] {
				continue
			}
			deg[to]--
			if deg[to] == 1 {
				queue = append(queue, to)
			}
		}
	}
	onCycle := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		if !removed[i] {
			onCycle[i] = true
		}
	}
	distA := make([]int, n+1)
	for i := range distA {
		distA[i] = -1
	}
	qa := []int{a}
	distA[a] = 0
	for h := 0; h < len(qa); h++ {
		v := qa[h]
		for _, to := range adj[v] {
			if distA[to] == -1 {
				distA[to] = distA[v] + 1
				qa = append(qa, to)
			}
		}
	}
	distB := make([]int, n+1)
	for i := range distB {
		distB[i] = -1
	}
	qb := []int{b}
	distB[b] = 0
	for h := 0; h < len(qb); h++ {
		v := qb[h]
		for _, to := range adj[v] {
			if distB[to] == -1 {
				distB[to] = distB[v] + 1
				qb = append(qb, to)
			}
		}
	}
	if distB[b] >= distA[b] {
		return "NO\n"
	}
	visited := make([]bool, n+1)
	q := []int{b}
	visited[b] = true
	escape := false
	for head := 0; head < len(q); head++ {
		v := q[head]
		if onCycle[v] {
			escape = true
			break
		}
		for _, to := range adj[v] {
			if !visited[to] && distB[to] < distA[to] {
				visited[to] = true
				q = append(q, to)
			}
		}
	}
	if escape {
		return "YES\n"
	}
	return "NO\n"
}

func runCase(exe, input, expect string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected %q got %q", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesH.txt")
	if err != nil {
		fmt.Println("could not read testcasesH.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		header := strings.Fields(scan.Text())
		if len(header) != 3 {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(header[0])
		a, _ := strconv.Atoi(header[1])
		b, _ := strconv.Atoi(header[2])
		edges := make([][2]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			parts := strings.Fields(scan.Text())
			if len(parts) != 2 {
				fmt.Println("bad file")
				os.Exit(1)
			}
			u, _ := strconv.Atoi(parts[0])
			v, _ := strconv.Atoi(parts[1])
			edges[i] = [2]int{u, v}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d %d %d\n", n, a, b)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", edges[i][0], edges[i][1])
		}
		exp := expected(n, a, b, edges)
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
