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

// Embedded testcases from testcasesE1.txt (input blocks separated by "---" and
// last line in each block was the expected answer). We compute expectations
// using the embedded solver, so the trailing expected lines are ignored.
const rawTestcases = `4 4 1
1 2
2 3
2 4
1 2
2 3
2 4
108
---
2 2 3
1 2
1 2
1 2
1 3
56
---
3 1 1
1 2
2 3
20
---
1 3 4
1 2
1 3
1 2
2 3
3 4
84
---
3 1 1
1 2
2 3
20
---
1 4 3
1 2
2 3
3 4
1 2
1 3
84
---
2 2 2
1 2
1 2
1 2
35
---
3 4 1
1 2
2 3
1 2
2 3
3 4
84
---
2 3 4
1 2
1 2
2 3
1 2
2 3
1 4
120
---
2 2 1
1 2
1 2
20
---
1 1 2
1 2
10
---
1 1 4
1 2
1 3
1 4
32
---
4 3 4
1 2
2 3
1 4
1 2
1 3
1 2
2 3
1 4
220
---
2 1 3
1 2
1 2
2 3
35
---
2 3 4
1 2
1 2
1 3
1 2
1 3
3 4
120
---
1 1 1
4
---
2 1 4
1 2
1 2
1 3
1 4
52
---
1 2 2
1 2
1 2
20
---
2 1 1
1 2
10
---
1 3 1
1 2
1 3
20
---
3 3 4
1 2
1 3
1 2
1 3
1 2
2 3
1 4
165
---
3 3 4
1 2
1 3
1 2
1 3
1 2
2 3
3 4
165
---
3 1 4
1 2
1 3
1 2
2 3
3 4
84
---
3 3 4
1 2
1 3
1 2
2 3
1 2
2 3
3 4
165
---
1 3 2
1 2
2 3
1 2
35
---
3 3 2
1 2
2 3
1 2
1 3
1 2
84
---
2 3 2
1 2
1 2
2 3
1 2
56
---
4 1 4
1 2
1 3
1 4
1 2
1 3
2 4
114
---
2 4 4
1 2
1 2
2 3
2 4
1 2
2 3
1 4
158
---
4 2 1
1 2
1 3
2 4
1 2
56
---
3 1 2
1 2
1 3
1 2
35
---
2 1 2
1 2
1 2
20
---
2 1 4
1 2
1 2
1 3
2 4
56
---
4 1 3
1 2
2 3
1 4
1 2
2 3
84
---
1 1 1
4
---
2 3 3
1 2
1 2
2 3
1 2
2 3
84
---
3 2 2
1 2
2 3
1 2
1 2
56
---
2 3 3
1 2
1 2
1 3
1 2
1 3
84
---
1 3 2
1 2
2 3
1 2
35
---
4 2 3
1 2
2 3
3 4
1 2
1 2
2 3
120
---
2 1 3
1 2
1 2
2 3
35
---
4 1 1
1 2
2 3
2 4
32
---
3 1 4
1 2
2 3
1 2
1 3
2 4
84
---
3 2 2
1 2
1 3
1 2
1 2
56
---
2 2 1
1 2
1 2
20
---
1 4 3
1 2
2 3
3 4
1 2
2 3
84
---
1 3 2
1 2
1 3
1 2
35
---
2 3 1
1 2
1 2
2 3
35
---
2 1 4
1 2
1 2
1 3
3 4
56
---
1 2 2
1 2
1 2
20
---
4 4 3
1 2
1 3
3 4
1 2
1 3
3 4
1 2
2 3
220
---
4 3 4
1 2
1 3
3 4
1 2
1 3
1 2
2 3
2 4
212
---
1 2 3
1 2
1 2
2 3
35
---
4 1 1
1 2
1 3
1 4
32
---
2 1 3
1 2
1 2
2 3
35
---
2 2 4
1 2
1 2
1 2
1 3
3 4
84
---
1 4 2
1 2
2 3
2 4
1 2
52
---
2 1 1
1 2
10
---
3 3 3
1 2
2 3
1 2
2 3
1 2
2 3
120
---
4 1 2
1 2
1 3
2 4
1 2
56
---
1 1 4
1 2
1 3
2 4
35
---
1 2 3
1 2
1 2
2 3
35
---
4 4 1
1 2
2 3
1 4
1 2
1 3
3 4
120
---
2 3 1
1 2
1 2
2 3
35
---
4 1 1
1 2
2 3
3 4
35
---
1 3 2
1 2
2 3
1 2
35
---
2 1 4
1 2
1 2
1 3
2 4
56
---
3 1 1
1 2
2 3
20
---
3 1 2
1 2
1 3
1 2
35
---
3 1 2
1 2
2 3
1 2
35
---
1 4 1
1 2
2 3
1 4
35
---
1 4 4
1 2
1 3
2 4
1 2
2 3
1 4
120
---
3 1 1
1 2
1 3
20
---
3 3 2
1 2
1 3
1 2
1 3
1 2
84
---
2 1 2
1 2
1 2
20
---
3 3 1
1 2
1 3
1 2
2 3
56
---
1 4 3
1 2
1 3
1 4
1 2
2 3
79
---
3 4 3
1 2
2 3
1 2
1 3
1 4
1 2
1 3
158
---
4 2 2
1 2
2 3
3 4
1 2
1 2
84
---
2 3 3
1 2
1 2
1 3
1 2
2 3
84
---
1 1 2
1 2
10
---
3 4 2
1 2
2 3
1 2
1 3
1 4
1 2
114
---
3 1 4
1 2
1 3
1 2
2 3
2 4
79
---
2 3 3
1 2
1 2
1 3
1 2
2 3
84
---
1 2 4
1 2
1 2
2 3
2 4
52
---
2 2 4
1 2
1 2
1 2
2 3
3 4
84
---
1 4 2
1 2
1 3
3 4
1 2
56
---
4 1 2
1 2
1 3
3 4
1 2
56
---
2 4 4
1 2
1 2
1 3
3 4
1 2
2 3
3 4
165
---
1 1 3
1 2
1 3
20
---
3 4 4
1 2
1 3
1 2
2 3
2 4
1 2
1 3
1 4
204
---
1 3 3
1 2
1 3
1 2
1 3
56
---
1 4 3
1 2
1 3
3 4
1 2
1 3
84
---
4 4 3
1 2
2 3
2 4
1 2
1 3
1 4
1 2
2 3
204
---
4 4 2
1 2
2 3
3 4
1 2
2 3
1 4
1 2
165
---
4 3 4
1 2
2 3
2 4
1 2
2 3
1 2
1 3
1 4
204
---
3 4 4
1 2
1 3
1 2
2 3
1 4
1 2
1 3
1 4
212
---
1 3 3
1 2
1 3
1 2
1 3
56
---
1 4 4
1 2
2 3
1 4
1 2
2 3
3 4
120
---
4 2 4
1 2
2 3
1 4
1 2
1 2
2 3
2 4
158`

type treeInfo struct {
	n        int
	D0       int64
	sumDepth []int64
	dist     [][]int
	maxSum   int64
}

func buildInfo(adj [][]int) treeInfo {
	n := len(adj)
	// D0 via subtree sizes
	var D0 int64
	sub := make([]int, n)
	var dfs func(u, p int)
	dfs = func(u, p int) {
		sub[u] = 1
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			dfs(v, u)
			sub[u] += sub[v]
			D0 += int64(sub[v]) * int64(n-sub[v])
		}
	}
	dfs(0, -1)

	sumDepth := make([]int64, n)
	dist := make([][]int, n)
	var maxSum int64
	for s := 0; s < n; s++ {
		d := make([]int, n)
		for i := range d {
			d[i] = -1
		}
		q := []int{s}
		d[s] = 0
		for qi := 0; qi < len(q); qi++ {
			u := q[qi]
			for _, v := range adj[u] {
				if d[v] == -1 {
					d[v] = d[u] + 1
					q = append(q, v)
				}
			}
		}
		dist[s] = d
		var sd int64
		for _, dd := range d {
			sd += int64(dd)
		}
		sumDepth[s] = sd
		if sd > maxSum {
			maxSum = sd
		}
	}

	return treeInfo{n: n, D0: D0, sumDepth: sumDepth, dist: dist, maxSum: maxSum}
}

func solveCase(input string) (string, error) {
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(input)))
	sc.Split(bufio.ScanWords)
	ints := make([]int, 0, 1024)
	for sc.Scan() {
		v, err := strconv.Atoi(sc.Text())
		if err != nil {
			return "", err
		}
		ints = append(ints, v)
	}
	if len(ints) < 3 {
		return "", fmt.Errorf("not enough data")
	}
	ns := []int{ints[0], ints[1], ints[2]}
	idx := 3
	adjList := make([][][]int, 3)
	for t := 0; t < 3; t++ {
		n := ns[t]
		if n < 1 {
			return "", fmt.Errorf("invalid n")
		}
		adj := make([][]int, n)
		for i := 0; i < n-1; i++ {
			if idx+1 >= len(ints) {
				return "", fmt.Errorf("truncated edges")
			}
			u, v := ints[idx]-1, ints[idx+1]-1
			idx += 2
			adj[u] = append(adj[u], v)
			adj[v] = append(adj[v], u)
		}
		adjList[t] = adj
	}

	info := make([]treeInfo, 3)
	for t := 0; t < 3; t++ {
		info[t] = buildInfo(adjList[t])
	}

	var totalD0 int64
	for _, in := range info {
		totalD0 += in.D0
	}

	var best int64
	for mid := 0; mid < 3; mid++ {
		ends := []int{}
		for t := 0; t < 3; t++ {
			if t != mid {
				ends = append(ends, t)
			}
		}
		a, k := ends[0], ends[1]
		ia, ij, ik := info[a], info[mid], info[k]
		sa, sj, sk := int64(ia.n), int64(ij.n), int64(ik.n)

		var bestInner int64
		for v := 0; v < ij.n; v++ {
			for w := 0; w < ij.n; w++ {
				val := sa*ij.sumDepth[v] + sk*ij.sumDepth[w] + sa*sk*int64(ij.dist[v][w])
				if val > bestInner {
					bestInner = val
				}
			}
		}

		cross := ia.maxSum*(sj+sk) + ik.maxSum*(sj+sa) + bestInner + sa*sj + sj*sk + 2*sa*sk
		if cross > best {
			best = cross
		}
	}

	ans := totalD0 + best
	return fmt.Sprint(ans), nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseTestcases() []string {
	blocks := strings.Split(strings.TrimSpace(rawTestcases), "\n---\n")
	inputs := make([]string, 0, len(blocks))
	for _, b := range blocks {
		lines := strings.Split(strings.TrimSpace(b), "\n")
		if len(lines) == 0 {
			continue
		}
		if len(lines) == 1 {
			inputs = append(inputs, lines[0])
			continue
		}
		// drop the last line (original expected answer) and keep the rest as input
		inputs = append(inputs, strings.Join(lines[:len(lines)-1], "\n"))
	}
	return inputs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := parseTestcases()
	for idx, input := range tests {
		expected, err := solveCase(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle error: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, strings.TrimSpace(input)+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
