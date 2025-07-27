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

type testCaseE struct {
	n, m    int
	t       []int
	queries [][2]int
}

func parseTestcases(path string) ([]testCaseE, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	in := bufio.NewReader(f)
	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return nil, err
	}
	cases := make([]testCaseE, T)
	for i := 0; i < T; i++ {
		var n, m int
		fmt.Fscan(in, &n, &m)
		t := make([]int, n+1)
		for j := 1; j <= n; j++ {
			fmt.Fscan(in, &t[j])
		}
		queries := make([][2]int, m-1)
		for j := 0; j < m-1; j++ {
			fmt.Fscan(in, &queries[j][0], &queries[j][1])
		}
		cases[i] = testCaseE{n: n, m: m, t: t, queries: queries}
	}
	return cases, nil
}

type DSU struct {
	parent []int
	size   []int
}

func NewDSU(n int) *DSU {
	d := &DSU{parent: make([]int, n+1), size: make([]int, n+1)}
	for i := 1; i <= n; i++ {
		d.parent[i] = i
		d.size[i] = 1
	}
	return d
}

func (d *DSU) Find(x int) int {
	for d.parent[x] != x {
		d.parent[x] = d.parent[d.parent[x]]
		x = d.parent[x]
	}
	return x
}

func solveCase(tc testCaseE) string {
	n, m := tc.n, tc.m
	t := append([]int(nil), tc.t...)
	items := make([][]int, m+1)
	for i := 1; i <= n; i++ {
		items[t[i]] = append(items[t[i]], i)
	}
	dsu := NewDSU(m)
	for i := 1; i <= m; i++ {
		dsu.size[i] = len(items[i])
	}
	ans := 0
	for i := 1; i < n; i++ {
		if t[i] != t[i+1] {
			ans++
		}
	}
	res := make([]int, 0, m)
	res = append(res, ans)
	for _, q := range tc.queries {
		a := dsu.Find(q[0])
		b := dsu.Find(q[1])
		if a != b {
			if dsu.size[a] < dsu.size[b] {
				a, b = b, a
			}
			for _, x := range items[b] {
				if x > 1 && dsu.Find(t[x-1]) == a {
					ans--
				}
				if x < n && dsu.Find(t[x+1]) == a {
					ans--
				}
			}
			dsu.parent[b] = a
			dsu.size[a] += dsu.size[b]
			items[a] = append(items[a], items[b]...)
			items[b] = nil
		}
		res = append(res, ans)
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func run(bin, input string) (string, string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return out.String(), errBuf.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcases("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for i := 1; i <= tc.n; i++ {
			if i > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.t[i]))
		}
		sb.WriteByte('\n')
		for _, q := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
		outStr, errStr, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", idx+1, err, errStr)
			os.Exit(1)
		}
		expected := solveCase(tc)
		if strings.TrimSpace(outStr) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, strings.TrimSpace(outStr))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
