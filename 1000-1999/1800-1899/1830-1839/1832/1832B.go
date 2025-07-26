package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		sort.Ints(arr)

		prefix := make([]int64, n+1)
		for i := 0; i < n; i++ {
			prefix[i+1] = prefix[i] + int64(arr[i])
		}

		var best int64
		for i := 0; i <= k; i++ {
			left := 2 * i
			right := n - (k - i)
			sum := prefix[right] - prefix[left]
			if sum > best {
				best = sum
			}
		}
		fmt.Fprintln(out, best)
	}
}
