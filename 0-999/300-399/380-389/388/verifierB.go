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

func parseGraph(out string) ([][]bool, error) {
	scan := bufio.NewScanner(strings.NewReader(out))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing n")
	}
	n, err := strconv.Atoi(strings.TrimSpace(scan.Text()))
	if err != nil {
		return nil, fmt.Errorf("bad n: %v", err)
	}
	if n < 2 || n > 1000 {
		return nil, fmt.Errorf("n out of bounds: %d", n)
	}
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("missing row %d", i)
		}
		line := strings.TrimSpace(scan.Text())
		if len(line) != n {
			return nil, fmt.Errorf("row %d length mismatch", i)
		}
		row := make([]bool, n)
		for j, ch := range line {
			if ch == 'Y' {
				row[j] = true
			} else if ch == 'N' {
				row[j] = false
			} else {
				return nil, fmt.Errorf("invalid char in row %d", i)
			}
		}
		adj[i] = row
	}
	if scan.Scan() {
		extra := strings.TrimSpace(scan.Text())
		if extra != "" {
			return nil, fmt.Errorf("extra output")
		}
	}
	return adj, nil
}

func countPaths(adj [][]bool) (int, int) {
	n := len(adj)
	dist := make([]int, n)
	for i := range dist {
		dist[i] = -1
	}
	cnt := make([]int, n)
	q := make([]int, 0, n)
	dist[0] = 0
	cnt[0] = 1
	q = append(q, 0)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		for u := 0; u < n; u++ {
			if !adj[v][u] {
				continue
			}
			if dist[u] == -1 {
				dist[u] = dist[v] + 1
				cnt[u] = cnt[v]
				q = append(q, u)
			} else if dist[u] == dist[v]+1 {
				cnt[u] += cnt[v]
			}
		}
	}
	return dist[1], cnt[1]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
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
		k, err := strconv.Atoi(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "bad testcase on line %d: %v\n", idx, err)
			os.Exit(1)
		}
		input := line + "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		adj, err := parseGraph(out.String())
		if err != nil {
			fmt.Printf("test %d: invalid output: %v\n", idx, err)
			os.Exit(1)
		}
		// validate graph is undirected simple
		n := len(adj)
		for i := 0; i < n; i++ {
			if adj[i][i] {
				fmt.Printf("test %d: self loop at %d\n", idx, i+1)
				os.Exit(1)
			}
			for j := i + 1; j < n; j++ {
				if adj[i][j] != adj[j][i] {
					fmt.Printf("test %d: asymmetry between %d and %d\n", idx, i+1, j+1)
					os.Exit(1)
				}
			}
		}
		dist, cnt := countPaths(adj)
		if dist == -1 {
			fmt.Printf("test %d: vertices 1 and 2 are disconnected\n", idx)
			os.Exit(1)
		}
		if cnt != k {
			fmt.Printf("test %d failed: expected %d shortest paths got %d\n", idx, k, cnt)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
