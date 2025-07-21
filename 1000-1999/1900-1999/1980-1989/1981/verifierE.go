package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
)

type edge struct {
	u, v int
	w    int
}

func intersect(l1, r1, l2, r2 int) bool {
	if l1 > r2 || l2 > r1 {
		return false
	}
	return true
}

func find(parent []int, x int) int {
	if parent[x] != x {
		parent[x] = find(parent, parent[x])
	}
	return parent[x]
}

func unite(parent []int, a, b int) bool {
	ra := find(parent, a)
	rb := find(parent, b)
	if ra == rb {
		return false
	}
	parent[rb] = ra
	return true
}

func mstWeight(n int, segs [][3]int) (int, bool) {
	var edges []edge
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if intersect(segs[i][0], segs[i][1], segs[j][0], segs[j][1]) {
				w := segs[i][2] - segs[j][2]
				if w < 0 {
					w = -w
				}
				edges = append(edges, edge{i, j, w})
			}
		}
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i].w < edges[j].w })
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
	}
	total := 0
	cnt := 0
	for _, e := range edges {
		if unite(parent, e.u, e.v) {
			total += e.w
			cnt++
			if cnt == n-1 {
				return total, true
			}
		}
	}
	return 0, false
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
	t, _ := strconv.Atoi(scan.Text())
	cases := make([][][3]int, t)
	answers := make([]int, t)
	okCase := make([]bool, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		segs := make([][3]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			l, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			r, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			a, _ := strconv.Atoi(scan.Text())
			segs[j] = [3]int{l, r, a}
		}
		if w, ok := mstWeight(n, segs); ok {
			answers[i] = w
			okCase[i] = true
		} else {
			okCase[i] = false
		}
		cases[i] = segs
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
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
		got, err := strconv.Atoi(outScan.Text())
		if err != nil {
			fmt.Printf("bad output for test %d\n", i+1)
			os.Exit(1)
		}
		if okCase[i] {
			if got != answers[i] {
				fmt.Printf("test %d failed: expected %d got %d\n", i+1, answers[i], got)
				os.Exit(1)
			}
		} else {
			if got != -1 {
				fmt.Printf("test %d failed: expected -1\n", i+1)
				os.Exit(1)
			}
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
