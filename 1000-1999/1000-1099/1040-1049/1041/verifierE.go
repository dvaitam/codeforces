package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type DSU struct{ parent, size []int }

func newDSU(n int) *DSU {
	p := make([]int, n)
	s := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
		s[i] = 1
	}
	return &DSU{p, s}
}

func (d *DSU) find(x int) int {
	if d.parent[x] != x {
		d.parent[x] = d.find(d.parent[x])
	}
	return d.parent[x]
}

func (d *DSU) union(a, b int) {
	a = d.find(a)
	b = d.find(b)
	if a == b {
		return
	}
	if d.size[a] < d.size[b] {
		a, b = b, a
	}
	d.parent[b] = a
	d.size[a] += d.size[b]
}

func solveE(n int, pairs [][2]int) (bool, [][2]int) {
	cnt := make([]int, n+1)
	for _, p := range pairs {
		cnt[p[0]]++
		cnt[p[1]]++
	}
	if cnt[n] != n-1 {
		return false, nil
	}
	for i := 1; i <= n; i++ {
		if cnt[i] > i {
			return false, nil
		}
	}
	zeros := make([]int, 0)
	for i := 1; i <= n; i++ {
		if cnt[i] == 0 {
			zeros = append(zeros, i)
		}
	}
	sort.Sort(sort.Reverse(sort.IntSlice(zeros)))
	sol := make([][2]int, 0)
	prv := n
	idx := 0
	cntCopy := append([]int(nil), cnt...)
	for i := n - 1; i >= 1; i-- {
		if cntCopy[i] == 0 {
			continue
		}
		for cntCopy[i] > 1 {
			if idx >= len(zeros) {
				return false, nil
			}
			u := prv
			v := zeros[idx]
			if v > i {
				return false, nil
			}
			sol = append(sol, [2]int{u, v})
			prv = v
			idx++
			cntCopy[i]--
		}
		sol = append(sol, [2]int{i, prv})
		prv = i
	}
	return true, sol
}

func computePairs(n int, edges [][2]int) [][2]int {
	res := make([][2]int, len(edges))
	for i, e := range edges {
		d := newDSU(n + 1)
		for j, f := range edges {
			if i == j {
				continue
			}
			d.union(f[0], f[1])
		}
		m1, m2 := 0, 0
		r1 := d.find(e[0])
		for v := 1; v <= n; v++ {
			if d.find(v) == r1 {
				if v > m1 {
					m1 = v
				}
			} else {
				if v > m2 {
					m2 = v
				}
			}
		}
		if m1 < m2 {
			res[i] = [2]int{m1, m2}
		} else {
			res[i] = [2]int{m2, m1}
		}
	}
	sort.Slice(res, func(i, j int) bool {
		if res[i][0] == res[j][0] {
			return res[i][1] < res[j][1]
		}
		return res[i][0] < res[j][0]
	})
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func parseEdges(out string) (bool, [][2]int, error) {
	parts := strings.Fields(out)
	if len(parts) == 0 {
		return false, nil, fmt.Errorf("empty output")
	}
	if parts[0] == "NO" {
		return false, nil, nil
	}
	if parts[0] != "YES" {
		return false, nil, fmt.Errorf("unexpected output")
	}
	if (len(parts)-1)%2 != 0 {
		return false, nil, fmt.Errorf("bad edges")
	}
	edges := make([][2]int, (len(parts)-1)/2)
	idx := 1
	for i := range edges {
		if idx+1 >= len(parts) {
			return false, nil, fmt.Errorf("bad edges")
		}
		var a, b int
		fmt.Sscan(parts[idx], &a)
		fmt.Sscan(parts[idx+1], &b)
		edges[i] = [2]int{a, b}
		idx += 2
	}
	return true, edges, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(5)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		pairs := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			a := rand.Intn(n-1) + 1
			b := rand.Intn(n-1) + 1
			if a > b {
				a, b = b, a
			}
			if a == b {
				b++
				if b > n {
					b = a
					a--
				}
			}
			if a < 1 {
				a = 1
			}
			if b > n {
				b = n
			}
			if a >= b {
				b = a + 1
				if b > n {
					b = a
				}
			}
			pairs[i] = [2]int{a, b}
		}
		input := fmt.Sprintf("%d\n", n)
		for i := 0; i < n-1; i++ {
			input += fmt.Sprintf("%d %d\n", pairs[i][0], pairs[i][1])
		}
		ok, expectedEdges := solveE(n, pairs)
		expectStr := "NO"
		if ok {
			expectStr = "YES"
			for _, e := range expectedEdges {
				expectStr += fmt.Sprintf(" %d %d", e[0], e[1])
			}
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			os.Exit(1)
		}
		gotOK, edges, err2 := parseEdges(out)
		if err2 != nil {
			fmt.Printf("test %d failed: %v\n", t, err2)
			os.Exit(1)
		}
		if ok != gotOK {
			fmt.Printf("test %d failed: expected %v got %v\n", t, ok, gotOK)
			os.Exit(1)
		}
		if ok {
			expPairs := computePairs(n, expectedEdges)
			gotPairs := computePairs(n, edges)
			if len(gotPairs) != len(expPairs) {
				fmt.Printf("test %d failed: wrong number of edges\n", t)
				os.Exit(1)
			}
			match := true
			for i := range expPairs {
				if expPairs[i] != gotPairs[i] {
					match = false
					break
				}
			}
			if !match {
				fmt.Printf("test %d failed: edges do not correspond\n", t)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
