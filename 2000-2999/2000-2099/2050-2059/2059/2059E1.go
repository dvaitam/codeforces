package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	out := bufio.NewWriter(os.Stdout)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		total := n * m

		flatA := make([]int, total)
		for i := 0; i < total; i++ {
			fmt.Fscan(in, &flatA[i])
		}

		bRows := make([][]int, n)
		for i := 0; i < n; i++ {
			row := make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(in, &row[j])
			}
			bRows[i] = row
		}

		var ops int64
		left, right := 0, total // current active segment of flatA is [left, right)

		for i := 0; i < n; i++ {
			row := bRows[i]

			bestK := m // default: fully replace the row
			if left < right {
				val := flatA[left]
				pos := -1
				for idx, x := range row {
					if x == val {
						pos = idx
						break
					}
				}
				if pos != -1 && left+m-pos <= right {
					match := true
					for j := 0; j < m-pos; j++ {
						if flatA[left+j] != row[pos+j] {
							match = false
							break
						}
					}
					if match {
						bestK = pos
					}
				}
			}

			ops += int64(bestK)
			left += m - bestK // discard matched prefix of the current segment
			right -= bestK    // drop consumed elements from the back
		}

		if t > 0 {
			fmt.Fprintln(out, ops)
		} else {
			fmt.Fprint(out, ops)
		}
	}

	out.Flush()
}
