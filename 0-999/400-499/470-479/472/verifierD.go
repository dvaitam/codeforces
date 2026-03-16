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
const testcasesDRaw = `
100
2
0 6
7 0
5
0 8 17 4 13
8 0 9 12 21
17 9 0 21 30
4 13 21 0 9
13 21 30 9 0
2
0 9
9 0
1
0
1
0
5
0 5 5 12 19
5 0 10 7 14
5 10 0 17 24
12 7 17 0 7
19 14 24 8 0
2
0 2
2 0
4
0 5 10 14
5 0 5 9
10 5 0 15
14 9 14 0
5
0 6 5 8 14
6 0 11 14 20
5 11 0 3 9
8 14 3 0 12
15 20 9 12 0
2
0 5
5 0
4
0 2 4 5
2 0 2 3
4 2 0 5
5 3 5 0
4
0 2 1 8
2 0 3 6
1 4 0 9
8 6 9 0
1
0
1
0
5
0 4 9 14 5
4 0 5 10 1
9 6 0 5 6
14 10 5 0 11
5 1 6 11 0
4
0 8 17 9
8 0 9 17
17 9 0 26
9 17 26 0
2
0 7
8 0
4
0 1 4 6
1 0 3 7
4 3 0 10
6 7 10 0
3
0 9 1
8 0 9
1 9 0
3
0 5 8
5 0 3
8 3 0
3
0 5 12
5 0 7
12 7 0
1
0
2
0 9
9 0
3
0 6 8
6 0 13
7 13 0
3
0 4 7
4 0 3
7 3 0
2
0 5
5 0
1
0
2
0 3
4 0
5
0 3 8 13 19
3 0 5 10 16
8 5 0 6 11
13 10 5 0 6
19 16 11 6 0
5
0 1 4 1 10
1 0 3 2 11
4 3 0 5 14
1 2 5 0 9
10 11 14 9 0
2
0 8
9 0
1
0
3
0 4 1
4 0 5
2 5 0
4
0 1 3 9
1 0 2 10
3 2 0 12
9 10 12 0
1
0
4
0 2 5 11
2 0 3 9
5 3 0 12
11 9 12 0
1
0
2
0 10
9 0
2
0 8
9 0
3
0 5 6
5 0 1
6 1 0
1
0
2
0 2
2 0
4
0 2 8 6
2 0 6 8
8 6 0 14
6 8 14 0
4
0 5 10 6
5 0 5 11
10 5 0 16
6 11 16 0
2
0 8
7 0
4
0 1 2 9
1 0 1 8
2 1 0 10
9 8 9 0
4
0 1 4 10
1 0 5 10
4 5 0 5
9 10 5 0
2
0 3
4 0
4
0 9 11 18
9 0 2 8
11 2 0 6
17 8 6 0
2
0 7
7 0
1
0
4
0 1 3 2
1 0 2 3
4 2 0 5
2 3 5 0
1
0
5
0 2 9 9 12
2 0 11 12 14
9 11 0 18 3
9 11 18 0 21
12 14 3 21 0
5
0 5 12 14 16
5 0 7 9 11
12 7 0 16 18
14 9 16 0 2
16 11 18 2 0
2
0 3
3 0
5
0 3 6 2 10
3 0 3 5 13
6 3 0 8 16
2 5 8 0 8
10 13 16 9 0
5
0 3 7 9 11
3 0 4 12 8
7 4 0 16 4
9 12 16 0 20
11 8 5 20 0
4
0 4 5 12
4 0 2 8
5 1 0 9
12 8 9 0
2
0 9
8 0
4
0 8 8 10
8 0 16 18
8 16 0 2
11 18 2 0
4
0 4 4 8
4 0 8 4
4 9 0 12
8 4 12 0
2
0 5
5 0
1
0
4
0 2 2 7
2 0 4 5
2 4 0 9
7 5 9 0
4
0 1 6 8
1 0 7 7
6 7 0 14
8 7 14 0
1
0
5
0 7 9 10 11
7 0 16 2 18
9 16 0 18 2
9 2 18 0 20
11 18 2 20 0
5
0 2 8 7 10
2 0 6 5 8
8 6 0 11 3
7 5 11 0 13
10 8 2 13 0
5
0 8 9 14 11
8 0 1 6 3
9 1 0 5 4
14 6 5 0 9
11 3 4 9 0
4
0 2 6 13
2 0 8 16
6 8 0 7
13 15 7 0
4
0 3 2 5
3 0 5 8
2 5 0 3
5 8 4 0
1
0
4
0 3 8 13
3 0 11 16
8 11 0 5
13 16 5 0
5
0 7 14 16 22
6 0 8 10 16
14 8 0 2 8
16 10 2 0 10
22 16 8 10 0
3
0 2 3
2 0 5
3 6 0
4
0 3 7 16
3 0 4 13
7 4 0 9
16 13 9 0
2
0 5
5 0
3
0 8 15
8 0 7
15 7 0
5
0 8 16 4 2
8 0 8 13 10
16 8 0 20 18
4 12 20 0 6
2 10 18 6 0
2
0 4
4 0
5
0 1 3 8 4
1 0 2 7 3
3 2 0 5 5
8 7 5 0 10
4 3 5 10 0
1
0
1
0
5
0 9 8 9 14
9 0 17 18 23
8 17 0 1 6
9 18 1 0 7
14 23 6 7 0
3
0 2 2
3 0 4
2 4 0
1
0
4
0 8 10 7
8 0 2 15
10 2 0 17
7 15 17 0
2
0 7
7 0
1
0
4
0 8 9 16
8 0 1 8
9 1 0 9
16 9 9 0
4
0 5 11 14
5 0 6 9
11 6 0 15
14 9 15 0
4
0 1 5 4
2 0 6 3
5 6 0 9
4 3 9 0
2
0 9
9 0
1
0
5
0 8 11 1 9
8 0 3 9 17
11 3 0 12 20
1 9 12 0 8
9 17 20 8 0
1
0
1
0
1
0
1
0
`


type edge struct {
	to int
	w  int64
}

func solveCase(n int, d [][]int64) string {
	const INF = int64(1e18)
	// Basic validation
	for i := 0; i < n; i++ {
		if d[i][i] != 0 {
			return "NO"
		}
		for j := 0; j < n; j++ {
			if d[i][j] != d[j][i] {
				return "NO"
			}
			if i != j && d[i][j] == 0 {
				return "NO"
			}
		}
	}
	if n == 1 {
		return "YES"
	}
	used := make([]bool, n)
	key := make([]int64, n)
	parent := make([]int, n)
	for i := 0; i < n; i++ {
		key[i] = INF
		parent[i] = -1
	}
	key[0] = 0
	for it := 0; it < n; it++ {
		u := -1
		best := INF
		for i := 0; i < n; i++ {
			if !used[i] && key[i] < best {
				best = key[i]
				u = i
			}
		}
		if u == -1 {
			return "NO"
		}
		used[u] = true
		for v := 0; v < n; v++ {
			if !used[v] && d[u][v] < key[v] {
				key[v] = d[u][v]
				parent[v] = u
			}
		}
	}
	adj := make([][]edge, n)
	for v := 1; v < n; v++ {
		u := parent[v]
		if u < 0 || key[v] <= 0 {
			return "NO"
		}
		w := key[v]
		adj[u] = append(adj[u], edge{v, w})
		adj[v] = append(adj[v], edge{u, w})
	}
	dist := make([]int64, n)
	q := make([]int, n)
	for src := 0; src < n; src++ {
		for i := 0; i < n; i++ {
			dist[i] = -1
		}
		head, tail := 0, 0
		q[tail] = src
		tail++
		dist[src] = 0
		for head < tail {
			u := q[head]
			head++
			for _, e := range adj[u] {
				v := e.to
				if dist[v] < 0 {
					dist[v] = dist[u] + e.w
					q[tail] = v
					tail++
				}
			}
		}
		for j := 0; j < n; j++ {
			if dist[j] != d[src][j] {
				return "NO"
			}
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data := []byte(testcasesDRaw)
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	type test struct {
		input    string
		expected string
	}
	tests := make([]test, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		d := make([][]int64, n)
		var buf bytes.Buffer
		fmt.Fprintln(&buf, n)
		for i := 0; i < n; i++ {
			row := make([]int64, n)
			for j := 0; j < n; j++ {
				if !scan.Scan() {
					fmt.Println("bad test file")
					os.Exit(1)
				}
				val, _ := strconv.ParseInt(scan.Text(), 10, 64)
				row[j] = val
				if j > 0 {
					buf.WriteByte(' ')
				}
				buf.WriteString(strconv.FormatInt(val, 10))
			}
			buf.WriteByte('\n')
			d[i] = row
		}
		tests[caseIdx] = test{input: buf.String(), expected: solveCase(n, d)}
	}
	for i, tc := range tests {
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc.input)
		out, err := cmd.Output()
		if err != nil {
			fmt.Println("execution failed:", err)
			os.Exit(1)
		}
		fields := strings.Fields(string(out))
		if len(fields) == 0 {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		if fields[0] != tc.expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, tc.expected, fields[0])
			os.Exit(1)
		}
		if len(fields) > 1 {
			fmt.Println("extra output detected")
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
