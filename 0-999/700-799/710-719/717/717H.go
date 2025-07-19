package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var n, m int
var g [][2]int
var l []int
var a [][]int
var maxT int
var t []int
var p []int
var col []int

func read() bool {
	reader := bufio.NewReader(os.Stdin)
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return false
	}
	g = make([][2]int, m)
	for i := 0; i < m; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		x--
		y--
		g[i][0], g[i][1] = x, y
	}
	l = make([]int, n)
	a = make([][]int, n)
	maxT = 0
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &l[i])
		a[i] = make([]int, l[i])
		for j := 0; j < l[i]; j++ {
			fmt.Fscan(reader, &a[i][j])
			a[i][j]--
			if a[i][j] > maxT {
				maxT = a[i][j]
			}
		}
	}
	maxT++
	return true
}

func solve() {
	rand.Seed(time.Now().UnixNano())
	// prepare permutation
	p = make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	t = make([]int, maxT)
	col = make([]int, n)
	fr := make([]int, 0, 21)

	for {
		// shuffle p[1..n-1]
		for i := 1; i < n; i++ {
			j := rand.Intn(i)
			p[i], p[j] = p[j], p[i]
		}
		// reset t
		for i := 0; i < maxT; i++ {
			t[i] = -1
		}
		// assign col by random rules
		for _, idx := range p {
			// gather free and pos0/pos1
			fr = fr[:0]
			pos0, pos1 := -1, -1
			for _, x := range a[idx] {
				if t[x] == -1 {
					fr = append(fr, x)
				} else if t[x] == 0 {
					pos0 = x
				} else if t[x] == 1 {
					pos1 = x
				}
			}
			if len(fr) > 0 && pos0 == -1 {
				k := rand.Intn(len(fr))
				t[fr[k]] = 0
				pos0 = fr[k]
				fr[k], fr[len(fr)-1] = fr[len(fr)-1], fr[k]
				fr = fr[:len(fr)-1]
			}
			if len(fr) > 0 && pos1 == -1 {
				k := rand.Intn(len(fr))
				t[fr[k]] = 1
				pos1 = fr[k]
			}
			if pos0 != -1 && pos1 != -1 {
				if rand.Int()&1 == 1 {
					col[idx] = pos0
				} else {
					col[idx] = pos1
				}
			} else if pos0 != -1 {
				col[idx] = pos0
			} else if pos1 != -1 {
				col[idx] = pos1
			} else {
				// should not happen
				col[idx] = a[idx][0]
			}
		}
		// count edges
		cnt := 0
		for i := 0; i < m; i++ {
			u, v := g[i][0], g[i][1]
			if t[col[u]] != t[col[v]] {
				cnt++
			}
		}
		if cnt*2 >= m {
			writer := bufio.NewWriter(os.Stdout)
			defer writer.Flush()
			for i := 0; i < n; i++ {
				fmt.Fprintf(writer, "%d ", col[i]+1)
			}
			fmt.Fprintln(writer)
			for i := 0; i < maxT; i++ {
				if t[i] == -1 {
					fmt.Fprintf(writer, "1 ")
				} else {
					fmt.Fprintf(writer, "%d ", t[i]+1)
				}
			}
			fmt.Fprintln(writer)
			return
		}
	}
}

func main() {
	if !read() {
		return
	}
	solve()
}
