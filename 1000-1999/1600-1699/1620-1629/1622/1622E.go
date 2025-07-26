package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		x := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &x[i])
		}
		strs := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &strs[i])
		}

		maskPerQ := make([]int, m)
		ones := make([]int, m)
		for j := 0; j < m; j++ {
			mask := 0
			for i := 0; i < n; i++ {
				if strs[i][j] == '1' {
					mask |= 1 << i
					ones[j]++
				}
			}
			maskPerQ[j] = mask
		}

		size := 1 << n
		pop := make([]int, size)
		for i := 1; i < size; i++ {
			pop[i] = pop[i>>1] + (i & 1)
		}

		bestVal := int64(math.MinInt64)
		bestOrder := make([]int, m)

		for mask := 0; mask < size; mask++ {
			constSum := 0
			for i := 0; i < n; i++ {
				if mask&(1<<i) != 0 {
					constSum += x[i]
				} else {
					constSum -= x[i]
				}
			}
			buckets := make([][]int, 2*n+1)
			for j := 0; j < m; j++ {
				a := pop[maskPerQ[j]&mask]*2 - ones[j]
				buckets[a+n] = append(buckets[a+n], j)
			}
			order := make([]int, m)
			pos := 1
			val := int64(constSum)
			for a := 0; a <= 2*n; a++ {
				aval := a - n
				for _, idx := range buckets[a] {
					order[idx] = pos
					val -= int64(aval * pos)
					pos++
				}
			}
			if val > bestVal {
				bestVal = val
				copy(bestOrder, order)
			}
		}

		for i := 0; i < m; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, bestOrder[i])
		}
		out.WriteByte('\n')
	}
}
