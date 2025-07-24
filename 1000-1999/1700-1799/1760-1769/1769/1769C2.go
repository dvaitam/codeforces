package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxDay = 1000005

var dp [maxDay]int

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	used := make([]int, 0, 2*200000)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		used = used[:0]
		best := 0
		for _, x := range arr {
			old := dp[x]
			valx := dp[x]
			if dp[x-1]+1 > valx {
				valx = dp[x-1] + 1
			}
			if valx != dp[x] {
				dp[x] = valx
				used = append(used, x)
			}
			if valx > best {
				best = valx
			}
			valx1 := dp[x+1]
			if old+1 > valx1 {
				valx1 = old + 1
			}
			if valx1 != dp[x+1] {
				dp[x+1] = valx1
				used = append(used, x+1)
			}
			if valx1 > best {
				best = valx1
			}
		}
		fmt.Fprintln(out, best)
		for _, idx := range used {
			dp[idx] = 0
		}
	}
}
