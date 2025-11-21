package main

import (
	"bufio"
	"fmt"
	"os"
)

const maxSum = 400005

func buildSieve(limit int) []bool {
	comp := make([]bool, limit+1)
	for i := 2; i*i <= limit; i++ {
		if !comp[i] {
			for j := i * i; j <= limit; j += i {
				comp[j] = true
			}
		}
	}
	return comp
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	isComposite := buildSieve(maxSum)

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		if n <= 4 {
			fmt.Fprintln(out, -1)
			continue
		}

		res := make([]int, 0, n)
		lastOdd := n
		if lastOdd%2 == 0 {
			lastOdd--
		}
		for i := 1; i <= n; i += 2 {
			res = append(res, i)
		}

		evens := make([]int, 0, n/2)
		for i := 2; i <= n; i += 2 {
			evens = append(evens, i)
		}

		idx := -1
		for i, val := range evens {
			sum := lastOdd + val
			if sum <= maxSum && isComposite[sum] {
				idx = i
				break
			}
		}
		if idx == -1 {
			fmt.Fprintln(out, -1)
			continue
		}

		res = append(res, evens[idx])
		for i, val := range evens {
			if i == idx {
				continue
			}
			res = append(res, val)
		}

		for i, val := range res {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, val)
		}
		fmt.Fprintln(out)
	}
}
