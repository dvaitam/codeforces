package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

type edge struct{ u, v int }

func possible(edges []edge) bool {
	maxv := 0
	for _, e := range edges {
		if e.u > maxv {
			maxv = e.u
		}
		if e.v > maxv {
			maxv = e.v
		}
	}
	adj := make([][]int, maxv+1)
	deg := make([]int, maxv+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
		deg[e.u]++
		deg[e.v]++
	}
	visited := make([]bool, maxv+1)
	comps := 0
	oddCount := 0
	for v := 1; v <= maxv; v++ {
		if deg[v] > 0 && !visited[v] {
			comps++
			stack := []int{v}
			visited[v] = true
			for len(stack) > 0 {
				u := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if deg[u]%2 == 1 {
					oddCount++
				}
				for _, to := range adj[u] {
					if !visited[to] {
						visited[to] = true
						stack = append(stack, to)
					}
				}
			}
		}
	}
	if comps > 2 {
		return false
	}
	if oddCount > 4 {
		return false
	}
	if comps == 2 {
		// each component must have 0 or 2 odd vertices
		visited = make([]bool, maxv+1)
		for v := 1; v <= maxv; v++ {
			if deg[v] > 0 && !visited[v] {
				stack := []int{v}
				visited[v] = true
				localOdd := 0
				for len(stack) > 0 {
					u := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					if deg[u]%2 == 1 {
						localOdd++
					}
					for _, to := range adj[u] {
						if !visited[to] {
							visited[to] = true
							stack = append(stack, to)
						}
					}
				}
				if localOdd != 0 && localOdd != 2 {
					return false
				}
			}
		}
	}
	return true
}

func parseInts(scan *bufio.Scanner, n int) ([]int, bool) {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if !scan.Scan() {
			return nil, false
		}
		v, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, false
		}
		arr[i] = v
	}
	return arr, true
}

func checkTrail(order []int, edges []edge, used []bool) bool {
	if len(order) == 0 {
		return false
	}
	if order[0] < 1 || order[0] > len(edges) {
		return false
	}
	used[order[0]-1] = true
	// orientation can be either
	start := edges[order[0]-1]
	// choose orientation to make path valid
	// We'll attempt both
	if checkTrailOrient(order, edges, used, start.u, start.v) {
		return true
	}
	for i := 0; i < len(order); i++ {
		used[order[i]-1] = false
	}
	return checkTrailOrient(order, edges, used, start.v, start.u)
}

func checkTrailOrient(order []int, edges []edge, used []bool, start, end int) bool {
	cur := end
	used[order[0]-1] = true
	for i := 1; i < len(order); i++ {
		id := order[i]
		if id < 1 || id > len(edges) || used[id-1] {
			return false
		}
		e := edges[id-1]
		if e.u == cur {
			cur = e.v
		} else if e.v == cur {
			cur = e.u
		} else {
			return false
		}
		used[id-1] = true
	}
	return true
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
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	cases := make([][]edge, t)
	possibleAns := make([]bool, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		var m int
		fmt.Sscan(scan.Text(), &m)
		edges := make([]edge, m)
		for j := 0; j < m; j++ {
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &edges[j].u)
			if !scan.Scan() {
				fmt.Println("bad file")
				os.Exit(1)
			}
			fmt.Sscan(scan.Text(), &edges[j].v)
		}
		cases[i] = edges
		possibleAns[i] = possible(edges)
	}
	for idx, edges := range cases {
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", len(edges))
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e.u, e.v)
		}
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		tok := outScan.Text()
		if tok == "-1" {
			if possibleAns[idx] {
				fmt.Printf("case %d should have solution\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if !possibleAns[idx] {
			fmt.Printf("case %d expected -1\n", idx+1)
			os.Exit(1)
		}
		L1, err := strconv.Atoi(tok)
		if err != nil || L1 <= 0 {
			fmt.Printf("case %d invalid L1\n", idx+1)
			os.Exit(1)
		}
		seq1, ok := parseInts(outScan, L1)
		if !ok {
			fmt.Printf("case %d bad sequence1\n", idx+1)
			os.Exit(1)
		}
		if !outScan.Scan() {
			fmt.Printf("case %d missing L2\n", idx+1)
			os.Exit(1)
		}
		L2, err := strconv.Atoi(outScan.Text())
		if err != nil || L2 <= 0 {
			fmt.Printf("case %d invalid L2\n", idx+1)
			os.Exit(1)
		}
		seq2, ok := parseInts(outScan, L2)
		if !ok {
			fmt.Printf("case %d bad sequence2\n", idx+1)
			os.Exit(1)
		}
		used := make([]bool, len(edges))
		if !checkTrail(seq1, edges, used) {
			fmt.Printf("case %d invalid trail1\n", idx+1)
			os.Exit(1)
		}
		if !checkTrail(seq2, edges, used) {
			fmt.Printf("case %d invalid trail2\n", idx+1)
			os.Exit(1)
		}
		for i, u := range used {
			if !u {
				fmt.Printf("case %d edge %d unused\n", idx+1, i+1)
				os.Exit(1)
			}
		}
		if outScan.Scan() {
			fmt.Printf("case %d: extra output\n", idx+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
