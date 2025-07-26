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
		arr := make([]int, n)
		for i := range arr {
			fmt.Fscan(in, &arr[i])
		}
		var q int
		fmt.Fscan(in, &q)

		states := make([][]int, 0, n+1)
		first := append([]int(nil), arr...)
		states = append(states, first)
		for step := 1; step <= n; step++ {
			prev := states[step-1]
			freq := make([]int, n+1)
			for _, v := range prev {
				freq[v]++
			}
			cur := make([]int, n)
			same := true
			for i, v := range prev {
				cur[i] = freq[v]
				if cur[i] != prev[i] {
					same = false
				}
			}
			states = append(states, cur)
			if same {
				break
			}
		}
		maxStep := len(states) - 1
		for ; q > 0; q-- {
			var x int
			var k int
			fmt.Fscan(in, &x, &k)
			if k > maxStep {
				k = maxStep
			}
			fmt.Fprintln(out, states[k][x-1])
		}
	}
}
