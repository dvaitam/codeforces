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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		maxVal := 0
		for i := 0; i < n; i++ {
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}

		best := 0
		for minVal := 1; minVal <= maxVal; minVal++ {
			pos := make([]int, 0)
			for i := 0; i < n; i++ {
				if a[i] == minVal {
					pos = append(pos, i)
				}
			}
			if len(pos) == 0 {
				continue
			}
			appIndex := 0
			starts := make([]int, len(pos))
			for idx, p := range pos {
				appIndex++
				starts[idx] = p - (appIndex - 1)
			}
			for medianVal := minVal; medianVal <= maxVal; medianVal++ {
				cnt := 0
				for j := 0; j < len(pos); j++ {
					p := pos[j]
					if a[p] < medianVal {
						continue
					}
					if a[p] == medianVal {
						cnt++
						if cnt%2 == 1 {
							length := p - starts[j] + 1
							res := medianVal - minVal
							if res > best {
								best = res
							}
						}
					} else {
						length := p - starts[j] + 1
						res := medianVal - minVal
						if res > best {
							best = res
						}
					}
				}
			}
		}
		fmt.Fprintln(out, best)
	}
}

