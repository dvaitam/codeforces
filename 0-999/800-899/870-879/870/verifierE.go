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

const mod int64 = 1000000007

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n), size: make([]int, n)}
	for i := 0; i < n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.Find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) Union(a, b int) {
	a = d.Find(a)
	b = d.Find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func solveE(n int, pts [][2]int) int64 {
	xMap := map[int]int{}
	yMap := map[int]int{}
	xi, yi := 0, 0
	for _, p := range pts {
		if _, ok := xMap[p[0]]; !ok {
			xMap[p[0]] = xi
			xi++
		}
		if _, ok := yMap[p[1]]; !ok {
			yMap[p[1]] = yi
			yi++
		}
	}
	total := xi + yi
	dsu := NewDSU(total)
	edges := make([][2]int, n)
	for i, p := range pts {
		xID := xMap[p[0]]
		yID := yMap[p[1]] + xi
		dsu.Union(xID, yID)
		edges[i] = [2]int{xID, yID}
	}
	edgeCnt := make([]int, total)
	for _, e := range edges {
		root := dsu.Find(e[0])
		edgeCnt[root]++
	}
	vertCnt := make([]int, total)
	for i := 0; i < total; i++ {
		root := dsu.Find(i)
		vertCnt[root]++
	}
	pow2 := make([]int64, total+1)
	pow2[0] = 1
	for i := 1; i <= total; i++ {
		pow2[i] = pow2[i-1] * 2 % mod
	}
	seen := map[int]bool{}
	ans := int64(1)
	for i := 0; i < total; i++ {
		root := dsu.Find(i)
		if seen[root] {
			continue
		}
		seen[root] = true
		v := vertCnt[root]
		e := edgeCnt[root]
		if e == v-1 {
			ans = ans * ((pow2[v] - 1 + mod) % mod) % mod
		} else {
			ans = ans * pow2[v] % mod
		}
	}
	return ans
}

func runCase(exe, input, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		pts := make([][2]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			x, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			y, _ := strconv.Atoi(scan.Text())
			pts[i] = [2]int{x, y}
		}
		inputBuilder := &strings.Builder{}
		fmt.Fprintf(inputBuilder, "%d\n", n)
		for _, p := range pts {
			fmt.Fprintf(inputBuilder, "%d %d\n", p[0], p[1])
		}
		expect := fmt.Sprintf("%d\n", solveE(n, pts))
		if err := runCase(exe, inputBuilder.String(), expect); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
