package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func bitsCount(x int) int {
	cnt := 0
	for x > 0 {
		x &= x - 1
		cnt++
	}
	return cnt
}

func solveCaseD(n, k, l int, xs []int, a []int) int {
	b := make([]bool, n+2)
	for _, v := range xs {
		if v >= 1 && v <= n {
			b[v] = true
		}
	}
	D := make([]int, 0)
	prev := false
	for i := 1; i <= n; i++ {
		if b[i] != prev {
			D = append(D, i)
		}
		prev = b[i]
	}
	if prev {
		D = append(D, n+1)
	}
	m := len(D)
	if m == 0 {
		return 0
	}
	if m%2 == 1 {
		return -1
	}
	N := n + 2
	adj := make([][]int, N)
	for s := 1; s <= n+1; s++ {
		for _, ai := range a {
			t := s + ai
			if t <= n+1 {
				adj[s] = append(adj[s], t)
				adj[t] = append(adj[t], s)
			}
		}
	}
	const INF = int(1e9)
	dist := make([][]int, m)
	q := make([]int, N)
	for i := 0; i < m; i++ {
		d := make([]int, N)
		for j := range d {
			d[j] = -1
		}
		qi, qj := 0, 0
		start := D[i]
		d[start] = 0
		q[qj] = start
		qj++
		for qi < qj {
			u := q[qi]
			qi++
			for _, v := range adj[u] {
				if d[v] == -1 {
					d[v] = d[u] + 1
					q[qj] = v
					qj++
				}
			}
		}
		dist[i] = make([]int, m)
		for j := 0; j < m; j++ {
			dv := d[D[j]]
			if dv >= 0 {
				dist[i][j] = dv
			} else {
				dist[i][j] = INF
			}
		}
	}
	maxMask := 1 << m
	dp := make([]int, maxMask)
	for i := 1; i < maxMask; i++ {
		dp[i] = INF
	}
	for mask := 1; mask < maxMask; mask++ {
		if bitsCount(mask)%2 == 1 {
			continue
		}
		var i0 int
		for bit := 0; bit < m; bit++ {
			if (mask>>bit)&1 == 1 {
				i0 = bit
				break
			}
		}
		for j := i0 + 1; j < m; j++ {
			if (mask>>j)&1 == 1 {
				m2 := mask ^ (1 << i0) ^ (1 << j)
				cost := dp[m2] + dist[i0][j]
				if cost < dp[mask] {
					dp[mask] = cost
				}
			}
		}
	}
	ans := dp[maxMask-1]
	if ans >= INF/2 {
		return -1
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	T, _ := strconv.Atoi(scan.Text())
	type one struct {
		n, k, l int
		xs      []int
		a       []int
	}
	cases := make([]one, T)
	expected := make([]int, T)
	for tc := 0; tc < T; tc++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		l, _ := strconv.Atoi(scan.Text())
		xs := make([]int, k)
		for i := 0; i < k; i++ {
			scan.Scan()
			xs[i], _ = strconv.Atoi(scan.Text())
		}
		a := make([]int, l)
		for i := 0; i < l; i++ {
			scan.Scan()
			a[i], _ = strconv.Atoi(scan.Text())
		}
		cases[tc] = one{n: n, k: k, l: l, xs: xs, a: a}
		expected[tc] = solveCaseD(n, k, l, xs, a)
	}
	for i, c := range cases {
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d %d\n", c.n, c.k, c.l)
		for _, x := range c.xs {
			fmt.Fprintf(&buf, "%d ", x)
		}
		buf.WriteByte('\n')
		for j, v := range c.a {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", v)
		}
		buf.WriteByte('\n')
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got, _ := strconv.Atoi(outScan.Text())
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d\n", i+1, expected[i], got)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected on case %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
