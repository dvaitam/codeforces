package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var x int64
		fmt.Fscan(reader, &n, &x)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		type pair struct {
			val int64
			cnt int64
		}
		q := make([]pair, n)
		for i := 0; i < n; i++ {
			q[i] = pair{a[i], 1}
		}
		sum := int64(0)
		idx := 0
		for idx < len(q) {
			p := q[idx]
			idx++
			sum += p.val * p.cnt
			if p.val%x != 0 {
				break
			}
			q = append(q, pair{p.val / x, p.cnt * x})
		}
		for idx < len(q) {
			p := q[idx]
			idx++
			sum += p.val * p.cnt
		}
		fmt.Fprintln(writer, sum)
	}
}
