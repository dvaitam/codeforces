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

func solveCase(n, m int, edges [][2]int) string {
	adj := make([][]int, n+1)
	for _, e := range edges {
		a, b := e[0], e[1]
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	vis := make([]bool, n+1)
	q := make([]int, 0)
	for i := 1; i <= n; i++ {
		if vis[i] {
			continue
		}
		q = append(q[:0], i)
		vis[i] = true
		nodes := 0
		edgesCnt := 0
		for len(q) > 0 {
			v := q[0]
			q = q[1:]
			nodes++
			edgesCnt += len(adj[v])
			for _, u := range adj[v] {
				if !vis[u] {
					vis[u] = true
					q = append(q, u)
				}
			}
		}
		if edgesCnt/2 != nodes*(nodes-1)/2 {
			return "NO"
		}
	}
	return "YES"
}

func runCandidate(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution failed: %v", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("missing n m")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			b, _ := strconv.Atoi(scan.Text())
			edges[i] = [2]int{a, b}
		}
		expected := solveCase(n, m, edges)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		out, err := runCandidate(os.Args[1], []byte(sb.String()))
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
		if out != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseIdx+1, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
