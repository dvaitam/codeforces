package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

// Embedded copy of testcasesE.txt so the verifier is self-contained.
const testcasesRaw = `100
1
1
3
6
2 6 6 3 3 5
7 20 2 19 6 14
7
7 6 7 5 3 5 4
17 9 2 1 12 15 11
7
4 5 2 5 2 2 2
1 6 11 6 5 17 17
6
5 6 5 2 4 4
17 12 19 12 12 15
3
2 3 3
15 17 8
8
5 8 6 8 8 6 8 8
8 11 6 20 9 16 10 10
7
3 6 2 4 5 3 6
20 3 11 1 7 4 2
1
1
19
4
1 2 3 2
7 2 14 2
1
1
12
3
1 3 1
3 4 3
1
1
1
6
3 2 2 6 2 5
1 13 19 2 8 5
1
1
12
2
2 2
16 1
5
4 5 5 1 3
13 20 5 16 8
2
2 1
1 15
3
3 3 2
16 17 11
3
2 2 2
20 14 1
3
3 1 2
2 5 6
3
1 2 3
8 17 2
4
2 4 1 3
3 19 8 20
6
3 6 4 3 5 1
5 2 13 14 6 4
2
1 1
4 1
3
1 1 1
1 17 15
8
5 7 4 4 7 7 1 1
14 17 19 6 4 16 12 1
2
2 2
12 10
1
1
4
2
2 1
1 15
1
1
16
8
4 2 1 5 1 6 5 2
8 16 7 4 19 12 13 15
3
2 2 1
9 4 4
2
2 2
7 4
1
1
2
8
5 6 8 3 6 5 8 8
14 16 10 13 8 6 16 20
5
5 4 1 5 5
4 3 12 6 18
3
2 1 1
2 5 10
7
2 6 6 6 3 4 2
17 10 4 5 18 14 4
6
5 2 6 5 3 2
6 15 8 13 12 19
3
2 2 3
1 20 13
3
2 3 1
16 9 13
5
4 4 3 5 3
3 8 18 20 7
7
7 6 4 6 1 3 4
17 15 6 4 1 13 7
7
2 1 4 7 5 7 7
7 9 19 19 7 16 20
3
1 3 3
14 16 9
3
2 3 1
3 12 1
8
2 8 6 8 5 8 1 2
20 12 6 13 9 5 2 6
8
7 8 5 3 1 5 8 1
12 2 18 13 19 15 7 10
8
3 8 5 2 5 6 5 6
10 13 17 3 17 7 13 20
3
3 3 1
10 2 8
8
4 5 1 2 2 7 6 4
11 12 3 11 15 12 6 16
8
5 8 3 8 4 5 6 3
4 8 16 7 12 6 12 5
3
1 2 3
13 13 11
5
5 5 5 3 4
10 18 20 3 12
5
4 4 2 3 3
15 16 3 6 11
7
2 1 1 3 2 3 1
14 1 18 11 8 20 13
5
4 2 3 3 2
16 4 5 7 11
5
2 4 3 3 1
11 7 8 8 20
1
1
12
1
1
6
2
2 2
9 5
6
5 5 1 3 6 6
20 13 8 2 13 16
8
6 2 8 7 8 3 7 7
17 15 2 4 15 19 5 4
3
1 2 2
15 1 9
2
2 1
6 1
3
2 3 1
11 15 2
8
4 2 8 3 1 3 1 1
7 18 1 17 11 17 8 5
6
4 1 2 5 1 2
4 15 7 2 20 7
7
3 5 6 7 4 6 5
17 6 17 4 5 7 6
7
2 3 3 4 2 4 2
13 11 10 4 18 4 16
5
3 5 4 3 2
14 5 18 4 1
4
2 2 4 1
5 1 9 16
1
1
5
6
1 6 2 1 2 6
18 6 3 15 10 7
3
2 3 2
17 19 3
7
4 6 6 1 4 3 6
4 9 1 7 14 11 9
7
5 5 6 2 4 7 2
6 15 15 12 13 16 20
5
5 2 5 4 4
7 16 19 11 10
2
1 2
20 16
4
2 3 2 3
4 1 1 7
6
1 3 5 3 6 6
11 15 3 14 16 1
5
5 5 2 2 2
6 20 13 3 19
8
5 2 8 8 4 3 5 4
7 20 11 19 20 13 17 14
4
2 1 3 2
5 20 13 14
2
2 2
13 16
7
3 2 2 2 1 5 5
3 20 18 1 2 13 14
7
2 5 3 1 3 5 3
17 16 19 3 15 8 9
1
1
16
1
1
5
4
3 2 1 2
10 4 18 18
2
1 2
5 2
5
5 3 4 1 5
12 11 4 20 12
2
2 2
9 16
5
5 5 2 1 1
11 14 1 12 18
1
1
18
7
4 4 2 7 2 2 5
2 1 19 12 6 10 1
1
1
19
4
4 1 3 1
20 3 8 8
4
1 1 4 1
17 9 19 8
4
1 2 2 1
1 4 5 3
5 5 1 1
4 1 1 2
2
4 4 5 5
5 3 2 5
4
3 2 1 2
4 2 2 5
5 3 2 1
5 1 3 2
2
4 3 2 2
1 1 1 1
2
3 3 3 3
5 4 2 2
2
2 5 1 5
5 5 4 1
1
5 5 5 5`

type frac struct {
	num int64
	den int64
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func reduce(a, b int64) frac {
	if a == 0 && b == 0 {
		return frac{0, 1}
	}
	if b == 0 {
		return frac{1, 0}
	}
	if a == 0 {
		return frac{0, 1}
	}
	g := gcd(abs(a), abs(b))
	a /= g
	b /= g
	if b < 0 {
		a = -a
		b = -b
	}
	return frac{a, b}
}

func computeExpected() (string, error) {
	in := bufio.NewReader(strings.NewReader(testcasesRaw))
	var out bytes.Buffer

	var n int
	fmt.Fscan(in, &n)
	A := make([]int64, n)
	B := make([]int64, n)
	C := make([]int64, n)
	D := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i], &B[i], &C[i], &D[i])
	}

	f1 := make([]frac, n)
	f2 := make([]frac, n)
	list := make([]frac, 0, 2*n)
	for i := 0; i < n; i++ {
		a1 := (A[i] + B[i]) * D[i]
		b1 := B[i] * C[i]
		f1[i] = reduce(a1, b1)
		list = append(list, f1[i])
		a2 := A[i] * D[i]
		b2 := B[i] * (C[i] + D[i])
		f2[i] = reduce(a2, b2)
		list = append(list, f2[i])
	}

	sort.Slice(list, func(i, j int) bool {
		if list[i].num != list[j].num {
			return list[i].num < list[j].num
		}
		return list[i].den < list[j].den
	})
	uniq := make([]frac, 0, len(list))
	for i, v := range list {
		if i == 0 || v != list[i-1] {
			uniq = append(uniq, v)
		}
	}
	idxMap := make(map[frac]int, len(uniq))
	for i, v := range uniq {
		idxMap[v] = i
	}

	type edge struct{ to, eid int }
	V := len(uniq)
	adj := make([][]edge, V)
	for i := 0; i < n; i++ {
		u := idxMap[f1[i]]
		v := idxMap[f2[i]]
		adj[u] = append(adj[u], edge{v, i})
		adj[v] = append(adj[v], edge{u, i})
	}

	vis := make([]bool, V)
	depth := make([]int, V)
	res := make([][2]int, 0, n/2)

	var dfs func(int, int, int) int
	dfs = func(u, pre, d int) int {
		vis[u] = true
		depth[u] = d
		cur := -1
		for _, e := range adj[u] {
			v := e.to
			tmp := -1
			if vis[v] {
				if depth[u] < depth[v] {
					tmp = e.eid
				}
			} else {
				tmp = dfs(v, e.eid, d+1)
			}
			if tmp == -1 {
				continue
			}
			if cur == -1 {
				cur = tmp
			} else {
				res = append(res, [2]int{cur, tmp})
				cur = -1
			}
		}
		if cur >= 0 && pre >= 0 {
			res = append(res, [2]int{cur, pre})
			cur = -1
			pre = -1
		}
		return pre
	}

	for i := 0; i < V; i++ {
		if !vis[i] {
			dfs(i, -1, 0)
		}
	}

	fmt.Fprintln(&out, len(res))
	for _, p := range res {
		fmt.Fprintf(&out, "%d %d\n", p[0]+1, p[1]+1)
	}

	return strings.TrimSpace(out.String()), nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	expected, err := computeExpected()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	got, err := runCandidate(bin, testcasesRaw)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if strings.TrimSpace(got) != expected {
		fmt.Printf("verification failed\nexpected:\n%s\n\ngot:\n%s\n", expected, got)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
