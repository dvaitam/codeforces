package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		var s string
		fmt.Fscan(reader, &s)
		l := make([]int, k)
		r := make([]int, k)
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &l[i])
		}
		for i := 0; i < k; i++ {
			fmt.Fscan(reader, &r[i])
		}

		// map each position to its segment index
		segIdx := make([]int, n+1)
		for idx := 0; idx < k; idx++ {
			for p := l[idx]; p <= r[idx]; p++ {
				segIdx[p] = idx
			}
		}

		pairCnt := make([]int, k)
		diffs := make([][]int, k)
		for i := 0; i < k; i++ {
			pc := (r[i] - l[i] + 1) / 2
			pairCnt[i] = pc
			diffs[i] = make([]int, pc+1)
		}

		var q int
		fmt.Fscan(reader, &q)
		for ; q > 0; q-- {
			var x int
			fmt.Fscan(reader, &x)
			seg := segIdx[x]
			pc := pairCnt[seg]
			if pc == 0 {
				continue
			}
			li, ri := l[seg], r[seg]
			y := li + ri - x
			a, b := x, y
			if a > b {
				a, b = b, a
			}
			start := a - li
			if start < ri-b {
				start = ri - b
			}
			end := b - li
			if end > ri-a {
				end = ri - a
			}
			if start < 0 {
				start = 0
			}
			if end >= pc {
				end = pc - 1
			}
			if start <= end {
				diffs[seg][start]++
				diffs[seg][end+1]--
			}
		}

		bytes := []byte(s)
		for seg := 0; seg < k; seg++ {
			pc := pairCnt[seg]
			if pc == 0 {
				continue
			}
			li, ri := l[seg], r[seg]
			cur := 0
			for i := 0; i < pc; i++ {
				cur += diffs[seg][i]
				if cur%2 != 0 {
					pos1 := li + i - 1
					pos2 := ri - i - 1
					bytes[pos1], bytes[pos2] = bytes[pos2], bytes[pos1]
				}
			}
		}

		fmt.Fprintln(writer, string(bytes))
	}
}
