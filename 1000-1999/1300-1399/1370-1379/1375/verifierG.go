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

func expectedAns(n int, edges [][2]int) int {
	adj := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	color := make([]int, n)
	for i := range color {
		color[i] = -1
	}
	q := []int{0}
	color[0] = 0
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, v := range adj[u] {
			if color[v] == -1 {
				color[v] = color[u] ^ 1
				q = append(q, v)
			}
		}
	}
	c0, c1 := 0, 0
	for _, c := range color {
		if c == 0 {
			c0++
		} else {
			c1++
		}
	}
	if c1 < c0 {
		c0 = c1
	}
	if c0 > 0 {
		c0--
	}
	return c0
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesG.txt")
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
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		expectEdges := (n - 1) * 2
		if len(parts) != 1+expectEdges {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, _ := strconv.Atoi(parts[1+2*i])
			v, _ := strconv.Atoi(parts[2+2*i])
			edges[i] = [2]int{u - 1, v - 1}
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
		}

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf, errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		out := strings.TrimSpace(outBuf.String())
		ansExp := expectedAns(n, edges)
		if out != fmt.Sprintf("%d", ansExp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, ansExp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
