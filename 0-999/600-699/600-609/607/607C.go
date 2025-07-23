package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	var s1, s2 string
	fmt.Fscan(in, &s1)
	fmt.Fscan(in, &s2)

	opp := map[byte]byte{'N': 'S', 'S': 'N', 'E': 'W', 'W': 'E'}
	type pair struct{ i, j int }
	vis := make(map[pair]struct{})
	q := make([]pair, 0)
	start := pair{0, 0}
	q = append(q, start)
	vis[start] = struct{}{}

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.i == n-1 && cur.j == n-1 {
			fmt.Fprintln(out, "YES")
			return
		}
		dirs := make(map[byte]struct{})
		if cur.i < n-1 {
			dirs[s1[cur.i]] = struct{}{}
		}
		if cur.i > 0 {
			dirs[opp[s1[cur.i-1]]] = struct{}{}
		}
		if cur.j < n-1 {
			dirs[s2[cur.j]] = struct{}{}
		}
		if cur.j > 0 {
			dirs[opp[s2[cur.j-1]]] = struct{}{}
		}
		for d := range dirs {
			ni, nj := cur.i, cur.j
			if cur.i < n-1 && s1[cur.i] == d {
				ni = cur.i + 1
			} else if cur.i > 0 && opp[s1[cur.i-1]] == d {
				ni = cur.i - 1
			}
			if cur.j < n-1 && s2[cur.j] == d {
				nj = cur.j + 1
			} else if cur.j > 0 && opp[s2[cur.j-1]] == d {
				nj = cur.j - 1
			}
			np := pair{ni, nj}
			if _, ok := vis[np]; !ok {
				vis[np] = struct{}{}
				q = append(q, np)
			}
		}
	}

	fmt.Fprintln(out, "NO")
}
