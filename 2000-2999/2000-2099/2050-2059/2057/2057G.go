package main

import (
	"bufio"
	"fmt"
	"os"
)

// A 5x5 periodic dominating pattern for the infinite grid.
var basePattern = [5][5]bool{
	{true, false, false, false, false},
	{false, false, true, false, false},
	{false, false, false, false, true},
	{false, true, false, false, false},
	{false, false, false, true, false},
}

type gridData struct {
	n, m int
	g    []byte
	s    int // number of free cells
	p    int // perimeter
}

func readGrid(in *bufio.Reader) gridData {
	var n, m int
	fmt.Fscan(in, &n, &m)
	g := make([]byte, n*m)
	s, p := 0, 0
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)
		for j := 0; j < m; j++ {
			ch := line[j]
			g[i*m+j] = ch
			if ch == '#' {
				s++
			}
		}
	}
	// compute perimeter in a separate pass
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g[i*m+j] != '#' {
				continue
			}
			if i == 0 || g[(i-1)*m+j] != '#' {
				p++
			}
			if i == n-1 || g[(i+1)*m+j] != '#' {
				p++
			}
			if j == 0 || g[i*m+j-1] != '#' {
				p++
			}
			if j == m-1 || g[i*m+j+1] != '#' {
				p++
			}
		}
	}
	return gridData{n: n, m: m, g: g, s: s, p: p}
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)

	patterns := make([][5][5]bool, 25)
	idx := 0
	for dr := 0; dr < 5; dr++ {
		for dc := 0; dc < 5; dc++ {
			for r := 0; r < 5; r++ {
				for c := 0; c < 5; c++ {
					patterns[idx][r][c] = basePattern[(r+dr)%5][(c+dc)%5]
				}
			}
			idx++
		}
	}

	for ; t > 0; t-- {
		data := readGrid(in)
		n, m, g := data.n, data.m, data.g
		if data.s == 0 {
			for i := 0; i < n; i++ {
				fmt.Fprintln(out, string(g[i*m:(i+1)*m]))
			}
			continue
		}

		colMod := make([]int, m)
		for j := 0; j < m; j++ {
			colMod[j] = j % 5
		}

		dom := make([]int, n*m)
		stamp := 1
		bestCnt := data.s + 1
		bestShift := 0

		mark := func(r, c int) {
			idx := r*m + c
			dom[idx] = stamp
			if r > 0 {
				dom[idx-m] = stamp
			}
			if r+1 < n {
				dom[idx+m] = stamp
			}
			if c > 0 {
				dom[idx-1] = stamp
			}
			if c+1 < m {
				dom[idx+1] = stamp
			}
		}

		for id, pat := range patterns {
			cnt := 0
			// initial pattern placement
			for i := 0; i < n; i++ {
				base := i * m
				rm := i % 5
				for j := 0; j < m; j++ {
					if g[base+j] == '#' && pat[rm][colMod[j]] {
						cnt++
						mark(i, j)
					}
				}
			}
			// fix uncovered cells
			for i := 0; i < n; i++ {
				base := i * m
				for j := 0; j < m; j++ {
					if g[base+j] == '#' && dom[base+j] != stamp {
						cnt++
						mark(i, j)
					}
				}
			}
			if cnt < bestCnt && cnt*5 <= data.s+data.p {
				bestCnt = cnt
				bestShift = id
			}
			stamp++
		}

		// build answer using best shift
		pat := patterns[bestShift]
		stamp++
		res := make([]byte, len(g))
		copy(res, g)
		mark = func(r, c int) {
			idx := r*m + c
			dom[idx] = stamp
			if r > 0 {
				dom[idx-m] = stamp
			}
			if r+1 < n {
				dom[idx+m] = stamp
			}
			if c > 0 {
				dom[idx-1] = stamp
			}
			if c+1 < m {
				dom[idx+1] = stamp
			}
		}

		for i := 0; i < n; i++ {
			base := i * m
			rm := i % 5
			for j := 0; j < m; j++ {
				if g[base+j] == '#' && pat[rm][colMod[j]] {
					res[base+j] = 'S'
					mark(i, j)
				}
			}
		}
		for i := 0; i < n; i++ {
			base := i * m
			for j := 0; j < m; j++ {
				if g[base+j] == '#' && dom[base+j] != stamp {
					res[base+j] = 'S'
					mark(i, j)
				}
			}
		}

		for i := 0; i < n; i++ {
			fmt.Fprintln(out, string(res[i*m:(i+1)*m]))
		}
	}
}
