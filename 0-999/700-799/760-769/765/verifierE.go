package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCase(n int, edges [][2]int) string {
	g := make([][]int, n)
	for _, e := range edges {
		u := e[0] - 1
		v := e[1] - 1
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	q := []int{0}
	dist[0] = 0
	for head := 0; head < len(q); head++ {
		v := q[head]
		for _, to := range g[v] {
			if dist[to] == -1 {
				dist[to] = dist[v] + 1
				q = append(q, to)
			}
		}
	}
	d := 0
	for i := 0; i < n; i++ {
		if len(g[i]) == 1 {
			if d == 0 {
				d = dist[i]
			} else {
				d = gcd(d, dist[i])
			}
		}
	}
	if d == 0 {
		return "0"
	}
	for d%2 == 0 {
		d /= 2
	}
	return fmt.Sprintf("%d", d)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &edges[i][0], &edges[i][1])
		}
		expected := solveCase(n, edges)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n-1; i++ {
			fmt.Fprintf(&input, "%d %d\n", edges[i][0], edges[i][1])
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
