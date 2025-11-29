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

const mod = 1000000007

const testcasesRaw = `8 1 1 2 0 2 3 0 3 4 1 3 5 0 5 6 0 5 7 0 4 8 1 
8 5 1 2 1 2 3 0 1 4 1 4 5 1 4 6 1 5 7 0 5 8 0 
3 2 1 2 0 2 3 0 
3 5 1 2 0 2 3 1 
7 5 1 2 1 2 3 1 1 4 1 4 5 0 4 6 1 4 7 1 
7 4 1 2 1 2 3 1 3 4 0 3 5 0 5 6 1 4 7 1 
4 5 1 2 1 1 3 1 3 4 1 
7 5 1 2 1 1 3 0 3 4 0 1 5 0 3 6 0 6 7 0 
8 5 1 2 1 1 3 0 1 4 1 1 5 0 3 6 1 2 7 0 6 8 0 
2 1 1 2 0 
2 1 1 2 1 
3 2 1 2 0 2 3 0 
8 2 1 2 0 1 3 1 3 4 0 3 5 1 4 6 0 3 7 1 5 8 0 
4 4 1 2 1 1 3 0 3 4 1 
8 1 1 2 1 1 3 1 2 4 1 2 5 1 3 6 1 5 7 1 6 8 0 
7 5 1 2 0 2 3 0 1 4 0 2 5 0 4 6 0 5 7 0 
3 2 1 2 0 2 3 0 
6 2 1 2 1 2 3 1 3 4 0 2 5 0 4 6 1 
3 1 1 2 0 1 3 0 
2 2 1 2 0 
7 2 1 2 1 1 3 1 1 4 1 4 5 0 3 6 0 6 7 0 
8 1 1 2 1 1 3 0 1 4 0 1 5 0 4 6 1 3 7 0 4 8 0 
7 4 1 2 1 2 3 1 1 4 0 1 5 0 5 6 1 1 7 1 
7 2 1 2 0 1 3 0 3 4 0 2 5 0 1 6 0 6 7 0 
2 3 1 2 0 
5 1 1 2 0 2 3 0 2 4 0 1 5 1 
2 1 1 2 0 
5 2 1 2 0 2 3 0 2 4 0 2 5 1 
3 1 1 2 0 1 3 1 
3 4 1 2 0 2 3 0 
8 1 1 2 1 1 3 0 3 4 0 3 5 0 5 6 0 6 7 0 6 8 0 
5 4 1 2 1 2 3 1 1 4 0 1 5 0 
7 4 1 2 0 2 3 0 2 4 0 4 5 1 4 6 1 5 7 1 
4 3 1 2 0 2 3 0 3 4 0 
6 3 1 2 1 1 3 1 1 4 1 4 5 1 5 6 0 
6 5 1 2 0 2 3 0 2 4 0 2 5 1 5 6 1 
8 2 1 2 1 2 3 0 2 4 0 1 5 0 2 6 0 3 7 1 5 8 0 
3 3 1 2 1 2 3 0 
7 2 1 2 0 2 3 0 3 4 1 4 5 0 3 6 0 2 7 1 
3 5 1 2 1 1 3 0 
7 4 1 2 0 1 3 1 2 4 0 1 5 0 1 6 1 4 7 0 
5 1 1 2 0 2 3 0 2 4 0 2 5 0 
6 4 1 2 1 1 3 0 2 4 0 3 5 0 4 6 0 
3 3 1 2 1 2 3 1 
4 2 1 2 1 1 3 0 3 4 1 
4 2 1 2 1 2 3 0 3 4 0 
8 3 1 2 1 2 3 1 3 4 0 4 5 1 5 6 1 4 7 1 4 8 1 
3 5 1 2 1 1 3 0 
5 4 1 2 0 1 3 0 1 4 1 3 5 1 
3 1 1 2 1 1 3 1 
5 5 1 2 0 2 3 0 2 4 0 2 5 0 
8 1 1 2 1 1 3 0 3 4 1 1 5 0 2 6 0 2 7 1 2 8 1 
5 3 1 2 0 2 3 0 2 4 1 1 5 1 
7 5 1 2 1 1 3 1 1 4 0 4 5 1 3 6 0 3 7 0 
4 1 1 2 0 2 3 0 3 4 1 
6 1 1 2 0 1 3 0 1 4 0 4 5 0 4 6 0 
6 3 1 2 1 2 3 1 3 4 0 3 5 0 2 6 1 
6 5 1 2 1 1 3 0 1 4 1 3 5 0 2 6 1 
6 3 1 2 0 2 3 0 2 4 0 3 5 0 2 6 1 
7 5 1 2 1 1 3 1 3 4 1 1 5 0 3 6 0 4 7 0 
5 4 1 2 0 1 3 0 3 4 0 2 5 0 
3 1 1 2 0 2 3 0 
2 4 1 2 1 
3 2 1 2 0 2 3 0 
5 3 1 2 0 2 3 0 1 4 0 1 5 0 
3 2 1 2 1 1 3 0 
6 1 1 2 1 2 3 1 1 4 1 4 5 0 5 6 0 
4 3 1 2 1 1 3 1 2 4 0 
4 1 1 2 1 2 3 0 2 4 0 
8 4 1 2 1 2 3 1 2 4 1 3 5 1 4 6 0 1 7 1 6 8 1 
6 4 1 2 0 1 3 1 3 4 0 2 5 0 5 6 0 
6 3 1 2 0 1 3 0 2 4 0 2 5 0 4 6 0 
3 1 1 2 0 2 3 1 
6 3 1 2 0 1 3 0 2 4 0 3 5 0 4 6 0 
7 4 1 2 1 1 3 0 1 4 0 4 5 0 1 6 0 1 7 0 
5 1 1 2 0 1 3 0 2 4 0 4 5 1 
6 4 1 2 0 2 3 0 3 4 1 2 5 1 4 6 1 
5 5 1 2 1 2 3 1 1 4 0 4 5 1 
7 3 1 2 1 1 3 0 2 4 0 1 5 0 1 6 0 2 7 0 
3 5 1 2 1 2 3 1 
5 5 1 2 0 2 3 0 2 4 0 2 5 1 
6 2 1 2 0 2 3 0 2 4 1 4 5 1 5 6 1 
8 1 1 2 1 1 3 1 1 4 1 4 5 1 3 6 0 3 7 0 4 8 1 
4 2 1 2 0 1 3 0 2 4 1 
6 1 1 2 1 2 3 1 1 4 1 4 5 1 1 6 0 
6 4 1 2 1 2 3 0 1 4 0 3 5 1 4 6 1 
4 5 1 2 0 2 3 1 3 4 1 
6 1 1 2 0 1 3 0 2 4 0 2 5 0 1 6 0 
5 4 1 2 0 2 3 0 2 4 0 1 5 0 
4 5 1 2 0 2 3 0 2 4 0 2 5 0 
5 3 1 2 0 1 3 0 1 4 0 3 5 1 
2 2 1 2 0 
7 4 1 2 0 2 3 0 1 4 0 1 5 0 4 6 1 4 7 0 
8 2 1 2 0 2 3 0 1 4 0 4 5 0 1 6 0 1 7 1 4 8 1 
8 3 1 2 1 2 3 1 1 4 0 2 5 0 5 6 1 4 7 1 1 8 1 
3 1 1 2 0 2 3 0 
7 3 1 2 0 2 3 0 1 4 0 3 5 0 2 6 0 5 7 0 
7 5 1 2 0 2 3 0 2 4 0 1 5 0 4 6 0 5 7 0 
8 4 1 2 0 2 3 0 4 4 0 3 5 0 4 6 0 3 7 0 7 8 0 
4 5 1 2 0 2 3 0 3 4 0 4 5 0 
2 4 1 2 0 
2 5 1 2 1 
4 5 1 2 1 2 3 0 2 4 0 
7 4 1 2 1 2 3 1 2 4 0 2 5 0 1 6 0 2 7 1 
4 5 1 2 1 2 3 1 1 4 0 
6 1 1 2 1 1 3 0 1 4 0 1 5 0 5 6 0 
8 2 1 2 1 1 3 1 3 4 0 3 5 0 4 6 0 4 7 0 5 8 0 
5 1 1 2 1 2 3 1 1 4 0 4 5 1 
5 2 1 2 0 2 3 0 3 4 0 2 5 0 
4 2 1 2 1 2 3 0 3 4 0 
6 5 1 2 1 1 3 1 3 4 1 1 5 0 4 6 0 
4 4 1 2 0 1 3 0 2 4 0 
6 2 1 2 1 2 3 0 1 4 0 1 5 0 5 6 1 
5 2 1 2 1 1 3 1 3 4 0 2 5 1 
5 4 1 2 1 2 3 0 3 4 0 4 5 1 
5 3 1 2 0 3 3 0 2 4 0 3 5 0 
7 1 1 2 0 1 3 0 3 4 0 4 5 0 5 6 0 6 7 1 
5 3 1 2 1 2 3 0 3 4 0 4 5 0 
5 4 1 2 1 2 3 1 3 4 0 2 5 0 
4 3 1 2 1 2 3 0 3 4 0 
8 4 1 2 0 1 3 0 2 4 0 4 5 0 5 6 0 6 7 0 7 8 0 
6 2 1 2 0 1 3 0 1 4 0 2 5 0 4 6 0 
2 1 1 2 0 
2 2 1 2 0 
7 4 1 2 1 2 3 0 1 4 0 1 5 1 1 6 0 6 7 1 
2 4 1 2 1 
5 1 1 2 1 2 3 1 3 4 1 1 5 0 
8 2 1 2 0 1 3 0 3 4 0 4 5 0 5 6 0 6 7 1 4 8 0 
7 2 1 2 1 2 3 0 1 4 0 1 5 0 5 6 1 2 7 1 
4 1 1 2 0 1 3 0 1 4 0 
6 3 1 2 1 1 3 1 2 4 0 1 5 0 1 6 0 
4 5 1 2 1 1 3 0 2 4 0 
4 1 1 2 0 2 3 0 2 4 0 
8 5 1 2 0 2 3 0 1 4 0 4 5 0 2 6 0 3 7 0 5 8 1 
7 1 1 2 0 2 3 0 3 4 0 1 5 1 2 6 1 1 7 1 
6 1 1 2 0 1 3 0 1 4 0 1 5 0 4 6 0 
5 1 1 2 1 2 3 0 2 4 0 2 5 1 
7 2 1 2 1 2 3 1 3 4 1 4 5 1 5 6 1 6 7 0 
3 2 1 2 1 1 3 0 
7 5 1 2 1 2 3 1 1 4 0 1 5 0 4 6 0 4 7 0 
2 1 1 2 0 
5 4 1 2 0 2 3 0 2 4 0 3 5 1 
4 5 1 2 0 2 3 0 3 4 0 
5 2 1 2 0 1 3 0 2 4 0 1 5 0 
5 4 1 2 1 2 3 1 1 4 0 2 5 0 
5 3 1 2 1 1 3 1 2 4 0 3 5 0 
6 3 1 2 1 1 3 1 3 4 0 1 5 0 2 6 0 
8 2 1 2 1 1 3 1 3 4 1 4 5 1 5 6 1 6 7 1 1 8 1 
2 3 1 2 0 
8 1 1 2 0 2 3 1 1 4 0 2 5 0 2 6 0 1 7 0 2 8 0 
4 1 1 2 0 1 3 0 2 4 0 
5 1 1 2 0 2 3 0 3 4 0 4 5 1 
8 5 1 2 0 2 3 0 2 4 0 4 5 1 1 6 0 1 7 1 2 8 1 
3 3 1 2 0 2 3 0 
7 3 1 2 0 1 3 0 1 4 0 4 5 1 1 6 1 2 7 0 
7 4 1 2 0 2 3 0 3 4 0 3 5 0 2 6 0 2 7 0 
3 5 1 2 0 1 3 1 
2 5 1 2 0 
3 5 1 2 0 1 3 0 
7 2 1 2 1 2 3 0 1 4 0 1 5 0 1 6 0 6 7 0 
7 4 1 2 1 2 3 1 2 4 0 1 5 0 2 6 0 3 7 1 
3 2 1 2 0 2 3 0 
5 2 1 2 0 2 3 0 1 4 0 4 5 1 
3 4 1 2 0 2 3 0 
5 4 1 2 1 1 3 0 1 4 0 4 5 0 
3 5 1 2 0 2 3 1 
4 3 1 2 0 2 3 0 1 4 0 
4 3 1 2 0 2 3 0 3 4 0 
3 4 1 2 0 1 3 0 
4 2 1 2 1 2 3 1 2 4 0`

type dsu struct {
	parent []int
	size   []int64
}

func newDSU(n int) *dsu {
	p := make([]int, n+1)
	sz := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &dsu{parent: p, size: sz}
}

func (d *dsu) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *dsu) unite(a, b int) {
	ra := d.find(a)
	rb := d.find(b)
	if ra == rb {
		return
	}
	d.parent[ra] = rb
	d.size[rb] += d.size[ra]
}

func modPow(a, b int64) int64 {
	var res int64 = 1
	a %= mod
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func referenceSolve(n, k int, edges [][3]int) int64 {
	ds := newDSU(n)
	for _, e := range edges {
		if e[2] == 0 {
			ds.unite(e[0], e[1])
		}
	}
	ans := modPow(int64(n), int64(k))
	for i := 1; i <= n; i++ {
		if ds.find(i) == i {
			ans = (ans - modPow(ds.size[i], int64(k)) + mod) % mod
		}
	}
	return ans
}

func runCase(bin, line string) error {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return fmt.Errorf("invalid test line")
	}
	n, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid n: %v", err)
	}
	k, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}
	edgeCount := (len(fields) - 2) / 3
	if edgeCount < n-1 {
		return fmt.Errorf("expected at least %d edges, got %d", n-1, edgeCount)
	}
	edges := make([][3]int, n-1)
	for i := 0; i < n-1; i++ {
		u, _ := strconv.Atoi(fields[2+3*i])
		v, _ := strconv.Atoi(fields[2+3*i+1])
		w, _ := strconv.Atoi(fields[2+3*i+2])
		edges[i] = [3]int{u, v, w}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, k)
	for _, e := range edges {
		fmt.Fprintf(&input, "%d %d %d\n", e[0], e[1], e[2])
	}

	expect := referenceSolve(n, k, edges)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse output %q", gotStr)
	}
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		if err := runCase(bin, line); err != nil {
			fmt.Printf("test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
