package main

import (
	"bufio"
	"fmt"
	"os"
)

const infVal = int(1e9)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		seq := make([]int, 0, n)
		for q := 1; q <= n; q++ {
			var k int
			fmt.Fscan(reader, &k)
			queries := make([]int, k)
			for i := 0; i < k; i++ {
				fmt.Fscan(reader, &queries[i])
			}
			// Build q-subsequence
			seq = seq[:0]
			for _, v := range p {
				if v <= q {
					seq = append(seq, v)
				}
			}
			m := len(seq)
			pre := make([]int, m)
			stack := make([]int, 0, m)
			for i := 0; i < m; i++ {
				for len(stack) > 0 && seq[stack[len(stack)-1]] <= seq[i] {
					stack = stack[:len(stack)-1]
				}
				if len(stack) == 0 {
					pre[i] = -infVal
				} else {
					pre[i] = stack[len(stack)-1]
				}
				stack = append(stack, i)
			}
			nxt := make([]int, m)
			stack = stack[:0]
			for i := m - 1; i >= 0; i-- {
				for len(stack) > 0 && seq[stack[len(stack)-1]] <= seq[i] {
					stack = stack[:len(stack)-1]
				}
				if len(stack) == 0 {
					nxt[i] = infVal
				} else {
					nxt[i] = stack[len(stack)-1]
				}
				stack = append(stack, i)
			}
			diff := make([]int, m)
			for i := 0; i < m; i++ {
				if pre[i] <= -infVal/2 || nxt[i] >= infVal/2 {
					diff[i] = infVal
				} else {
					diff[i] = nxt[i] - pre[i]
				}
			}
			for _, x := range queries {
				sum := 0
				for i := 0; i < m; i++ {
					d := diff[i]
					if d > x {
						sum += x
					} else {
						sum += d
					}
				}
				fmt.Fprintln(writer, sum)
			}
		}
	}
}
