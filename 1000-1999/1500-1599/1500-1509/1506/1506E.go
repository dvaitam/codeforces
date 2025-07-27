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
		q := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &q[i])
		}

		minAns := make([]int, n)
		maxAns := make([]int, n)
		usedMin := make([]bool, n+1)
		usedMax := make([]bool, n+1)

		// Build lexicographically minimal permutation
		nextMin := 1
		for i := 0; i < n; i++ {
			if i == 0 || q[i] != q[i-1] {
				minAns[i] = q[i]
				usedMin[q[i]] = true
			} else {
				for nextMin <= n && usedMin[nextMin] {
					nextMin++
				}
				minAns[i] = nextMin
				usedMin[nextMin] = true
			}
		}

		// Build lexicographically maximal permutation
		cur := 0
		stack := []int{}
		for i := 0; i < n; i++ {
			if i == 0 || q[i] != q[i-1] {
				maxAns[i] = q[i]
				// push numbers between previous prefix max and current - 1
				for x := cur + 1; x < q[i]; x++ {
					if !usedMax[x] {
						stack = append(stack, x)
						usedMax[x] = true
					}
				}
				usedMax[q[i]] = true
				cur = q[i]
			} else {
				// pop largest available number less than current prefix max
				idx := len(stack) - 1
				maxAns[i] = stack[idx]
				stack = stack[:idx]
			}
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, minAns[i])
		}
		out.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				out.WriteByte(' ')
			}
			fmt.Fprint(out, maxAns[i])
		}
		out.WriteByte('\n')
	}
}
