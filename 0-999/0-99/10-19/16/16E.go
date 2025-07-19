package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// Read probability matrix
	a := make([][]float64, n)
	for i := 0; i < n; i++ {
		a[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &a[i][j])
		}
	}
	// dp over subsets
	size := 1 << n
	d := make([]float64, size)
	ans := make([]float64, n)
	full := size - 1
	d[full] = 1.0
	// temporary slice for indices
	s := make([]int, n)
	// iterate masks from full down to 1
	for mask := full; mask > 0; mask-- {
		probCur := d[mask]
		if probCur == 0 {
			continue
		}
		// collect bits
		cnt := 0
		m := mask
		for i := 0; i < n; i++ {
			if m&(1<<i) != 0 {
				s[cnt] = i
				cnt++
			}
		}
		if cnt == 1 {
			ans[s[0]] = probCur
		} else {
			totalPairs := float64(cnt * (cnt - 1))
			fac := 2.0 / totalPairs
			// for each pair i<j
			for x := 0; x < cnt-1; x++ {
				i := s[x]
				for y := x + 1; y < cnt; y++ {
					j := s[y]
					// i loses to j: remove i
					d[mask^(1<<i)] += probCur * fac * a[j][i]
					// j loses to i: remove j
					d[mask^(1<<j)] += probCur * fac * a[i][j]
				}
			}
		}
	}
	// Output answers
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for i := 0; i < n; i++ {
		fmt.Fprintf(writer, "%.12f", ans[i])
		if i+1 < n {
			writer.WriteByte(' ')
		}
	}
	writer.WriteByte('\n')
}
