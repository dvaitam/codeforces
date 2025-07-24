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

	var k, q int
	if _, err := fmt.Fscan(reader, &k, &q); err != nil {
		return
	}

	const maxP = 1000
	ans := make([]int, maxP+1)
	prev := make([]float64, k+1)
	curr := make([]float64, k+1)
	prev[0] = 1.0
	idx := 1
	n := 0
	kf := float64(k)
	for idx <= maxP {
		n++
		curr[0] = 0
		for j := 1; j <= k; j++ {
			curr[j] = prev[j]*float64(j)/kf + prev[j-1]*float64(k-j+1)/kf
		}
		prev, curr = curr, prev
		p := prev[k]
		for idx <= maxP && p >= float64(idx)/2000.0-1e-12 {
			ans[idx] = n
			idx++
		}
	}

	for i := 0; i < q; i++ {
		var pi int
		fmt.Fscan(reader, &pi)
		if pi < 1 {
			pi = 1
		} else if pi > 1000 {
			pi = 1000
		}
		fmt.Fprintln(writer, ans[pi])
	}
}
