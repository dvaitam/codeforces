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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &r[i])
		}
		var m int
		fmt.Fscan(in, &m)
		b := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}
		bestR := 0
		sum := 0
		for _, v := range r {
			sum += v
			if sum > bestR {
				bestR = sum
			}
		}
		bestB := 0
		sum = 0
		for _, v := range b {
			sum += v
			if sum > bestB {
				bestB = sum
			}
		}
		fmt.Fprintln(out, bestR+bestB)
	}
}
