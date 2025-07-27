package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pair struct {
	val int
	idx int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		nm := n * m
		arr := make([]int, nm)
		for i := 0; i < nm; i++ {
			fmt.Fscan(in, &arr[i])
		}
		pairs := make([]Pair, nm)
		for i := 0; i < nm; i++ {
			pairs[i] = Pair{val: arr[i], idx: i}
		}
		sort.Slice(pairs, func(i, j int) bool {
			if pairs[i].val == pairs[j].val {
				return pairs[i].idx < pairs[j].idx
			}
			return pairs[i].val < pairs[j].val
		})

		posRow := make([]int, nm)
		posCol := make([]int, nm)

		for r := 0; r < n; r++ {
			slice := pairs[r*m : (r+1)*m]
			sort.Slice(slice, func(i, j int) bool {
				if slice[i].val == slice[j].val {
					return slice[i].idx > slice[j].idx
				}
				return slice[i].val < slice[j].val
			})
			for c := 0; c < m; c++ {
				idx := slice[c].idx
				posRow[idx] = r
				posCol[idx] = c
			}
		}

		occupied := make([][]bool, n)
		for i := 0; i < n; i++ {
			occupied[i] = make([]bool, m)
		}
		var total int
		for i := 0; i < nm; i++ {
			r := posRow[i]
			c := posCol[i]
			cnt := 0
			for j := 0; j < c; j++ {
				if occupied[r][j] {
					cnt++
				}
			}
			total += cnt
			occupied[r][c] = true
		}
		fmt.Fprintln(out, total)
	}
}
