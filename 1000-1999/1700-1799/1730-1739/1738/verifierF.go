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

func solveF(n int, edges [][2]int) []int {
	g := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	color := make([]int, n+1)
	cur := 0
	q := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		if color[i] != 0 {
			continue
		}
		cur++
		q = q[:0]
		q = append(q, i)
		color[i] = cur
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			for _, v := range g[u] {
				if color[v] == 0 {
					color[v] = cur
					q = append(q, v)
				}
			}
		}
	}
	return color[1:]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
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
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+2*m {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, m)
		for i := 0; i < m; i++ {
			u, _ := strconv.Atoi(parts[2+2*i])
			v, _ := strconv.Atoi(parts[3+2*i])
			edges[i] = [2]int{u, v}
		}
		expect := solveF(n, edges)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotParts := strings.Fields(strings.TrimSpace(out.String()))
		if len(gotParts) != n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n, len(gotParts))
			os.Exit(1)
		}
		got := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(gotParts[i])
			got[i] = v
		}
		for i := 0; i < n; i++ {
			if got[i] != expect[i] {
				fmt.Printf("test %d failed\nexpected %v\ngot %v\n", idx, expect, got)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
