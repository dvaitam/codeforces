package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, na, nb, nc int
		fmt.Fscan(in, &n, &na, &nb, &nc)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] > arr[j] })
		groups := [][]int{make([]int, 0, na), make([]int, 0, nb), make([]int, 0, nc)}
		sums := []int{0, 0, 0}
		caps := []int{na, nb, nc}
		for _, v := range arr {
			idx := -1
			minSum := int(^uint(0) >> 1)
			for j := 0; j < 3; j++ {
				if caps[j] > 0 && sums[j] < minSum {
					minSum = sums[j]
					idx = j
				}
			}
			if idx == -1 {
				break
			}
			groups[idx] = append(groups[idx], v)
			sums[idx] += v
			caps[idx]--
		}
		if caps[0] == 0 && caps[1] == 0 && caps[2] == 0 {
			sa, sb, sc := sums[0], sums[1], sums[2]
			if sa < sb+sc && sb < sa+sc && sc < sa+sb {
				fmt.Fprintln(out, "YES")
				for i := 0; i < 3; i++ {
					for j, x := range groups[i] {
						if j > 0 {
							fmt.Fprint(out, " ")
						}
						fmt.Fprint(out, x)
					}
					fmt.Fprintln(out)
				}
				continue
			}
		}
		fmt.Fprintln(out, "NO")
	}
}
