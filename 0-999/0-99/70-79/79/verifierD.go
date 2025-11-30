package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

// Embedded testcases (same format as original file).
const embeddedTestcases = `100
96 6 7
4 17 7 29 59 73
15 23 24 25 30 32 33
28 6 6
57 29 8 32 43 74
21 34 16 2 20 19
12 5 8
9 4 11 1 5
6 4 16 1 5 10 7 13
4 2 14
1 3
2 1
59 1 5
42
9 29 52 50 47
12 6 14
5 6 4 8 7 2
7 23 28 16 15 6 8 11 14 12 25 2 24 1
41 5 2
16 19 21 4 37
25 17
72 2 8
66 37
2 53 28 32 38 63 6 49
45 2 8
15 25
26 39 32 41 24 7 8 2
33 2 8
15 31
32 22 11 18 31 6 8 1
19 6 6
19 10 8 15 18 8
2 13 4 6 7 18
31 5 9
24 27 8 14 18
2 1 12 7 11 4 15 3 5
32 5 4
24 32 12 31 5
9 23 27 20
53 6 14
3 20 17 14 4 10
15 5 32 10 21 20 11 4 26 1 30 24 19 16
18 6 8
14 18 8 5 11 17
6 15 8 9 17 4 18 5
80 1 6
19
1 2 7 10 4 3
29 1 7
10
26 1 27 3 4 15 2 5 9 6 20 18 22
87 6 7
44 69 10 29 51 82
24 3 17 12 10 27 13
31 2 15
29 8
23 26 7 14 3 4 5 9 18 17 27 25 22 19 28
10 1 5
7
8 6 1 10 4
16 1 9
4
1 6 4 9 12 15 2 3 13
40 6 9
8 11 12 6 10 5
1 13 4 12 7 10 6 3 11
62 1 6
14
4 2 11 6 10 3
16 1 4
6
1 2 4 5
34 2 4
8 4
2 4 6 5
9 5 3
1 4 6 7 5
2 5 1
39 6 4
11 38 26 7 32 23
6 11 17 33
72 2 3
51 68
1 2 3
60 4 7
18 30 58 45
1 7 6 5 4 3 2
21 5 15
10 3 20 11 14
4 3 8 2 9 6 1 7 10 5 12 11 14 15 13
44 5 6
26 30 37 35 5
16 6 3 13 2 11
28 1 15
6
2 5 4 3 13 7 6 9 14 10 12 8 11 15 1
57 6 16
7 11 45 6 39 33
7 10 4 14 6 15 12 8 3 16 13 11 5 2 1 9
52 1 16
45
7 11 1 9 12 14 15 2 16 8 5 10 4 13 6 3
68 1 11
60
3 9 8 11 1 10 5 2 7 4 6
64 5 17
41 46 34 6 29
8 14 12 1 3 6 7 2 5 9 16 4 17 13 11 10 15
80 4 14
49 79 12 64
1 5 14 2 8 3 11 10 6 9 13 12 7 4
93 2 15
38 84
12 15 8 4 11 10 2 14 7 5 6 1 9 3 13
77 3 2
67 62 69
21 2
27 6 14
5 14 6 22 10 23
15 6 12 18 3 17 5 7 11 10 16 20 13 4
58 2 6
48 11
3 6 4 1 5 2
48 4 4
42 4 39 12
9 7 1 3
60 4 7
44 4 46 51
3 1 4 6 2 7 5
33 6 10
5 16 9 15 24 29
15 8 10 20 7 11 17 21 9 3
68 6 10
48 13 47 28 39 53
2 10 7 5 6 1 3 9 4 8
42 2 9
28 32
5 6 1 3 9 2 7 4 8
18 4 7
5 6 12 10
3 1 4 5 2 6 7
24 5 12
15 6 22 23 4
2 12 11 5 9 6 7 10 3 4 1 8
70 5 9
47 4 1 67 49
2 5 1 7 3 4 9 6 8
78 2 11
31 47
9 3 10 6 4 1 11 5 7 2
36 4 16
19 5 26 11
4 3 9 15 8 2 10 16 7 6 1 11 12 5 13 14
6 1 7
5
4 7 1 5 3 2 6
95 3 16
79 17 23
7 16 12 4 3 6 5 8 2 15 13 9 11 10 1 14
12 1 4
7
3 4 2 1
45 3 6
3 39 34
3 2 5 4 6 1
36 1 5
34
1 3 2 5 4
53 4 2
21 49 8 48
9 2
60 5 17
11 28 48 37 13
14 3 6 9 17 16 2 10 12 8 5 4 1 15
28 2 5
24 9
3 5 2 1 4
91 1 15
12
4 13 5 9 10 6 8 3 1 2 7 15 14 12 11
22 1 4
6
1 3 4 2
68 2 7
16 65
7 1 6 5 3 2 4
32 6 11
22 4 2 11 9 14
6 1 8 2 3 7 10 5 9 11 4
90 1 8
35
3 4 5 8 7 2 6 1
96 1 6
50
4 2 6 3 1 5
74 5 9
30 61 4 25 1
1 3 6 8 9 2 4 5 7
80 1 13
12
2 12 10 3 9 5 8 1 11 7 4 13 6
48 6 12
10 13 43 47 24 14
9 10 4 3 8 11 7 12 1 5 6
10 2 17
8 6
7 11 4 17 12 16 9 15 2 10 6 14 1 5 3 8 13
76 6 10
17 22 72 3 64 52
7 6 9 5 4 3 10 2 1 8
78 2 4
30 36
4 2 1 3
72 3 13
48 63 61
1 7 4 6 9 11 13 12 10 3 8 5 2
7 5 8
4 6 2 3 1
7 6 5 8 1 2 3 4
31 3 9
1 13 12
2 6 1 5 3 9 7 8 4
85 4 7
26 45 43 70
4 2 3 1 6 5 7
41 5 7
7 19 18 14 30
1 5 6 7 2 3 4
39 4 2
17 35 33 30
2 1
10 1 8
10
7 2 1 4 3 8 5 6
96 3 17
6 3 30
9 3 4 5 7 2 10 17 12 6 13 8 14 15 16 1 11
4 1 6
1
4 1 2 3 6 5
74 1 12
34
4 2 7 9 11 8 6 10 12 3 1 5
28 4 16
9 15 22 24
8 9 13 15 12 11 1 3 7 5 16 14 6 4 2 10
35 6 13
7 29 20 32 22 30
4 8 2 6 10 3 9 11 13 7 5 1 12
58 3 8
5 47 40
1 8 2 7 3 4 5 6
51 1 3
49
3 2 1
26 1 3
2
3 2 1
45 6 6
35 18 26 3 21 6
6 1 5 4 3 2
52 2 10
21 35
5 9 2 7 10 4 8 3 6 1
69 1 15
34
14 15 12 7 10 5 4 2 9 13 11 6 1 8 3
35 4 8
13 10 30 1
1 2 8 4 6 5 7 3
15 1 14
5
1 3 14 4 8 5 6 9 10 13 2 12 11 7
53 1 13
32
1 2 9 4 3 12 11 7 5 8 10 6 13
76 2 6
40 66
5 2 1 3 6 4`

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
	data := []byte(embeddedTestcases)
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
