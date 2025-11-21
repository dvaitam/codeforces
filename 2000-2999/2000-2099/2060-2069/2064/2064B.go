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
		a := make([]int, n)
		freq := make([]int, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			freq[a[i]]++
		}

		bestLen := 0
		bestL, bestR := -1, -1

		i := 0
		for i < n {
			if freq[a[i]] != 1 {
				i++
				continue
			}
			j := i
			for j < n && freq[a[j]] == 1 {
				j++
			}
			if j-i > bestLen {
				bestLen = j - i
				bestL = i
				bestR = j - 1
			}
			i = j
		}

		if bestLen == 0 {
			fmt.Fprintln(out, 0)
		} else {
			fmt.Fprintf(out, "%d %d\n", bestL+1, bestR+1)
		}
	}
}
