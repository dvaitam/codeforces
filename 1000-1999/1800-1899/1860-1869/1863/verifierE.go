package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseE struct {
	n     int
	m     int
	k     int64
	h     []int64
	edges [][2]int
}

func generateTestsE(num int) []testCaseE {
	rand.Seed(time.Now().UnixNano())
	tests := make([]testCaseE, num)
	for i := 0; i < num; i++ {
		n := rand.Intn(8) + 1
		m := rand.Intn(n*(n-1)/2 + 1)
		k := rand.Int63n(20) + 1
		h := make([]int64, n)
		for j := 0; j < n; j++ {
			h[j] = rand.Int63n(k)
		}
		edges := make([][2]int, 0, m)
		used := map[[2]int]bool{}
		for len(edges) < m {
			a := rand.Intn(n)
			b := rand.Intn(n)
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			e := [2]int{a, b}
			if used[e] {
				continue
			}
			used[e] = true
			edges = append(edges, e)
		}
		tests[i] = testCaseE{n: n, m: m, k: k, h: h, edges: edges}
	}
	return tests
}

func solveE(tc testCaseE) string {
	n := tc.n
	// tc.m is not needed in the reference implementation
	k := tc.k
	h := tc.h
	adj := make([][]int, n)
	indeg := make([]int, n)
	for _, e := range tc.edges {
		a := e[0]
		b := e[1]
		adj[a] = append(adj[a], b)
		indeg[b]++
	}
	q := make([]int, 0, n)
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	for idx := 0; idx < len(q); idx++ {
		u := q[idx]
		for _, v := range adj[u] {
			indeg[v]--
			if indeg[v] == 0 {
				q = append(q, v)
			}
		}
	}
	dp := make([]int64, n)
	for i := n - 1; i >= 0; i-- {
		u := q[i]
		for _, v := range adj[u] {
			w := (h[v] - h[u]) % k
			if w < 0 {
				w += k
			}
			if dp[u] < dp[v]+w {
				dp[u] = dp[v] + w
			}
		}
	}
	indeg = make([]int, n)
	for u := 0; u < n; u++ {
		for _, v := range adj[u] {
			indeg[v]++
		}
	}
	sources := []int{}
	for i := 0; i < n; i++ {
		if indeg[i] == 0 {
			sources = append(sources, i)
		}
	}
	s := len(sources)
	if s == 0 {
		return "0"
	}
	hs := make([]int64, s)
	ls := make([]int64, s)
	for i, idx := range sources {
		hs[i] = h[idx]
		ls[i] = dp[idx]
	}
	order := make([]int, s)
	for i := 0; i < s; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return hs[order[i]] < hs[order[j]] })
	hsorted := make([]int64, s)
	lsorted := make([]int64, s)
	for i, idx := range order {
		hsorted[i] = hs[idx]
		lsorted[i] = ls[idx]
	}
	prefMaxHL := make([]int64, s)
	prefMinH := make([]int64, s)
	var mx int64 = -1 << 60
	var mn int64 = 1 << 60
	for i := 0; i < s; i++ {
		val := hsorted[i] + lsorted[i]
		if val > mx {
			mx = val
		}
		prefMaxHL[i] = mx
		if hsorted[i] < mn {
			mn = hsorted[i]
		}
		prefMinH[i] = mn
	}
	suffMaxHL := make([]int64, s)
	suffMinH := make([]int64, s)
	mx = -1 << 60
	mn = 1 << 60
	for i := s - 1; i >= 0; i-- {
		val := hsorted[i] + lsorted[i]
		if val > mx {
			mx = val
		}
		suffMaxHL[i] = mx
		if hsorted[i] < mn {
			mn = hsorted[i]
		}
		suffMinH[i] = mn
	}
	unique := make(map[int64]bool)
	res := int64(1 << 62)
	for _, c := range hsorted {
		if unique[c] {
			continue
		}
		unique[c] = true
		idx := sort.Search(len(hsorted), func(i int) bool { return hsorted[i] >= c })
		max1 := int64(-1 << 60)
		min1 := int64(1 << 60)
		if idx < s {
			max1 = suffMaxHL[idx] - c
			min1 = suffMinH[idx] - c
		}
		max2 := int64(-1 << 60)
		min2 := int64(1 << 60)
		if idx > 0 {
			max2 = prefMaxHL[idx-1] - c + k
			min2 = prefMinH[idx-1] - c + k
		}
		maxv := max1
		if max2 > maxv {
			maxv = max2
		}
		minv := min1
		if min2 < minv {
			minv = min2
		}
		diff := maxv - minv
		if diff < res {
			res = diff
		}
	}
	return fmt.Sprint(res)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTestsE(100)
	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d %d\n", tc.n, tc.m, tc.k)
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.h[i])
		}
		input.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0]+1, e[1]+1)
		}
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("binary execution failed:", err)
		os.Exit(1)
	}
	outputs := strings.Split(strings.TrimSpace(string(out)), "\n")
	if len(outputs) != len(tests) {
		fmt.Printf("expected %d lines of output, got %d\n", len(tests), len(outputs))
		os.Exit(1)
	}
	for i, tc := range tests {
		expected := solveE(tc)
		if strings.TrimSpace(outputs[i]) != expected {
			fmt.Printf("mismatch on test %d: expected %s got %s\n", i+1, expected, outputs[i])
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
