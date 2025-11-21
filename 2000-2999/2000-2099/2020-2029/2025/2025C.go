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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Ints(arr)
		values := make([]int, 0)
		counts := make([]int, 0)
		for i := 0; i < n; {
			j := i
			for j < n && arr[j] == arr[i] {
				j++
			}
			values = append(values, arr[i])
			counts = append(counts, j-i)
			i = j
		}
		m := len(values)
		best := 0
		sum := 0
		l := 0
		for r := 0; r < m; r++ {
			if r > 0 && values[r]-values[r-1] > 1 {
				// reset window
				sum = 0
				l = r
			}
			sum += counts[r]
			for r-l+1 > k {
				sum -= counts[l]
				l++
			}
			if sum > best {
				best = sum
			}
		}
		fmt.Fprintln(out, best)
	}
}
