package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		// compute divisors of n excluding n
		divs := []int{}
		for d := 1; d*d <= n; d++ {
			if n%d == 0 {
				if d < n {
					divs = append(divs, d)
				}
				if n/d != d && n/d < n {
					divs = append(divs, n/d)
				}
			}
		}
		sort.Ints(divs)

		sums := make(map[int][]int64)
		maxSum := make(map[int]int64)

		for _, g := range divs {
			arr := make([]int64, g)
			for i := 0; i < n; i++ {
				arr[i%g] += a[i]
			}
			var m int64
			for _, v := range arr {
				if v > m {
					m = v
				}
			}
			sums[g] = arr
			maxSum[g] = m
		}

		// initial answer
		var ans int64
		for _, g := range divs {
			val := int64(g) * maxSum[g]
			if val > ans {
				ans = val
			}
		}
		fmt.Fprintln(writer, ans)

		for ; q > 0; q-- {
			var p int
			var x int64
			fmt.Fscan(reader, &p, &x)
			p--
			diff := x - a[p]
			if diff != 0 {
				a[p] = x
				for _, g := range divs {
					idx := p % g
					old := sums[g][idx]
					nv := old + diff
					sums[g][idx] = nv
					if old == maxSum[g] {
						if nv > old {
							maxSum[g] = nv
						} else {
							var m int64
							for _, v := range sums[g] {
								if v > m {
									m = v
								}
							}
							maxSum[g] = m
						}
					} else if nv > maxSum[g] {
						maxSum[g] = nv
					}
				}
			}

			ans = 0
			for _, g := range divs {
				val := int64(g) * maxSum[g]
				if val > ans {
					ans = val
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
