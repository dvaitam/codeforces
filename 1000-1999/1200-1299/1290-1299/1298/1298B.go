package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	seen := make(map[int]bool)
	res := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		if !seen[arr[i]] {
			seen[arr[i]] = true
			res = append(res, arr[i])
		}
	}

	// reverse res to restore order of rightmost occurrences
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, len(res))
	for i, v := range res {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, v)
	}
	fmt.Fprintln(out)
	out.Flush()
}
