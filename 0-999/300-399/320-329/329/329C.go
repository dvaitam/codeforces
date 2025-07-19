package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n, m    int
	b       [][]int
	cl      []int
	leftCnt []int
	bns     [][2]int
	ans     [][2]int
	flagOK  bool
)

func f(x, y, z int) bool {
	// check if any edge exists among x,y,z
	for _, v := range b[x] {
		if v == y || v == z {
			return true
		}
	}
	for _, v := range b[z] {
		if v == y {
			return true
		}
	}
	return false
}

// backtrack to pick lm edges among ln nodes
func btrk(ln, lm, loc, p, q int) {
	if loc == lm {
		flagOK = true
		return
	}
	if p == ln {
		return
	}
	if p == q {
		btrk(ln, lm, loc, p+1, 0)
		return
	}
	u := cl[p]
	v := cl[q]
	if u > 0 && v > 0 && u != v && leftCnt[u] < 2 && leftCnt[v] < 2 {
		// ensure no original edge between u and v
		ok := true
		for _, w := range b[u] {
			if w == v {
				ok = false
				break
			}
		}
		if ok {
			bns[loc][0] = u
			bns[loc][1] = v
			leftCnt[u]++
			leftCnt[v]++
			btrk(ln, lm, loc+1, p, q+1)
			leftCnt[u]--
			leftCnt[v]--
			if flagOK {
				return
			}
		}
	}
	btrk(ln, lm, loc, p, q+1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &n, &m)
	b = make([][]int, n+1)
	cl = make([]int, n)
	leftCnt = make([]int, n+1)
	bns = make([][2]int, m)
	ans = make([][2]int, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		b[u] = append(b[u], v)
		b[v] = append(b[v], u)
	}
	for i := 0; i < n; i++ {
		cl[i] = i + 1
	}
	if n <= 7 {
		btrk(n, m, 0, 0, 0)
		if !flagOK {
			fmt.Fprintln(writer, -1)
		} else {
			for i := 0; i < m; i++ {
				fmt.Fprintf(writer, "%d %d\n", bns[i][0], bns[i][1])
			}
		}
		return
	}
	var i int
	for i = 0; i < n-7; i += 3 {
		// find triple j,k,l
		var jj, kk, ll int
		found := false
		for j := 0; j < 7 && !found; j++ {
			for k := 0; k < j && !found; k++ {
				for l := 0; l < k; l++ {
					if !f(cl[n-i-1-j], cl[n-i-1-k], cl[n-i-1-l]) {
						jj, kk, ll = j, k, l
						found = true
						break
					}
				}
			}
		}
		// assign edges
		ans[i][0] = cl[n-i-1-jj]
		ans[i][1] = cl[n-i-1-kk]
		ans[i+1][0] = cl[n-i-1-jj]
		ans[i+1][1] = cl[n-i-1-ll]
		ans[i+2][0] = cl[n-i-1-kk]
		ans[i+2][1] = cl[n-i-1-ll]
		// remove used nodes
		cl[n-i-1-jj] = 0
		cl[n-i-1-kk] = 0
		cl[n-i-1-ll] = 0
		// compact cl array
		for t := 0; t < 3; t++ {
			idx := n - i - 1 - t
			if cl[idx] != 0 {
				// find position to move
				pos := n - i - 4
				for pos >= 0 && cl[pos] != 0 {
					pos--
				}
				if pos >= 0 {
					cl[pos] = cl[idx]
				}
			}
		}
	}
	if i >= m {
		for j := 0; j < m; j++ {
			fmt.Fprintf(writer, "%d %d\n", ans[j][0], ans[j][1])
		}
	} else {
		for j := 0; j < i; j++ {
			fmt.Fprintf(writer, "%d %d\n", ans[j][0], ans[j][1])
		}
		// backtrack for remaining
		btrk(n-i, m-i, 0, 0, 0)
		for j := 0; j < m-i; j++ {
			fmt.Fprintf(writer, "%d %d\n", bns[j][0], bns[j][1])
		}
	}
}
