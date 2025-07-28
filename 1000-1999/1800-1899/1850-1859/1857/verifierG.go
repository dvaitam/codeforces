package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const mod int64 = 998244353

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n)
	sz := make([]int, n)
	for i := range p {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{p, sz}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	ra, rb := d.find(a), d.find(b)
	if ra == rb {
		return
	}
	if d.size[ra] < d.size[rb] {
		ra, rb = rb, ra
	}
	d.parent[rb] = ra
	d.size[ra] += d.size[rb]
}

func powMod(a, e int64) int64 {
	a %= mod
	if a < 0 {
		a += mod
	}
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func expected(n int, S int64, edges [][3]int64) string {
	freq := make(map[int64]int64)
	for _, e := range edges {
		freq[e[2]]++
	}
	sort.Slice(edges, func(i, j int) bool { return edges[i][2] < edges[j][2] })
	dsu := NewDSU(n)
	counts := make(map[int64]int64)
	for _, e := range edges {
		u := dsu.find(int(e[0]) - 1)
		v := dsu.find(int(e[1]) - 1)
		if u != v {
			counts[e[2]] += int64(dsu.size[u]) * int64(dsu.size[v])
			dsu.union(u, v)
		}
	}
	ans := int64(1)
	zero := false
	for w, c := range counts {
		exp := c - freq[w]
		if exp <= 0 {
			continue
		}
		diff := S - w
		if diff <= 0 {
			zero = true
			break
		}
		ans = ans * powMod(diff%mod, exp) % mod
	}
	if zero {
		return "0"
	}
	return fmt.Sprintf("%d", ans%mod)
}

func runCase(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesG.txt")
	if err != nil {
		fmt.Println("failed to read testcasesG.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("case %d missing nS\n", caseNum)
			os.Exit(1)
		}
		header := strings.Fields(scan.Text())
		if len(header) != 2 {
			fmt.Printf("case %d malformed header\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(header[0])
		S, _ := strconv.ParseInt(header[1], 10, 64)
		edges := make([][3]int64, n-1)
		for i := 0; i < n-1; i++ {
			if !scan.Scan() {
				fmt.Printf("case %d missing edge %d\n", caseNum, i)
				os.Exit(1)
			}
			parts := strings.Fields(scan.Text())
			if len(parts) != 3 {
				fmt.Printf("case %d malformed edge\n", caseNum)
				os.Exit(1)
			}
			u, _ := strconv.ParseInt(parts[0], 10, 64)
			v, _ := strconv.ParseInt(parts[1], 10, 64)
			w, _ := strconv.ParseInt(parts[2], 10, 64)
			edges[i] = [3]int64{u, v, w}
		}
		inputLines := []string{"1", fmt.Sprintf("%d %d", n, S)}
		for _, e := range edges {
			inputLines = append(inputLines, fmt.Sprintf("%d %d %d", e[0], e[1], e[2]))
		}
		input := strings.Join(inputLines, "\n") + "\n"
		want := expected(n, S, edges)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
