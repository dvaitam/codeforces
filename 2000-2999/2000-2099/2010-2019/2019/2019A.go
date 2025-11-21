package main

import (
	"bufio"
	"fmt"
	"os"
)

const negInf = -1 << 30

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		uniq := make(map[int]struct{})
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
			uniq[arr[i]] = struct{}{}
		}
		bestScore := 0
		for m := range uniq {
			dp0 := [2]int{0, negInf}
			dp1 := [2]int{negInf, negInf}
			for _, val := range arr {
				new0 := [2]int{negInf, negInf}
				new1 := [2]int{negInf, negInf}
				for has := 0; has <= 1; has++ {
					best := dp0[has]
					if dp1[has] > best {
						best = dp1[has]
					}
					if best > new0[has] {
						new0[has] = best
					}
				}
				if val <= m {
					for has := 0; has <= 1; has++ {
						if dp0[has] <= negInf/2 {
							continue
						}
						flag := has
						if val == m {
							flag = 1
						}
						cand := dp0[has] + 1
						if cand > new1[flag] {
							new1[flag] = cand
						}
					}
				}
				dp0 = new0
				dp1 = new1
			}
			bestCount := dp0[1]
			if dp1[1] > bestCount {
				bestCount = dp1[1]
			}
			if bestCount > negInf/2 {
				score := m + bestCount
				if score > bestScore {
					bestScore = score
				}
			}
		}
		fmt.Fprintln(out, bestScore)
	}
}
