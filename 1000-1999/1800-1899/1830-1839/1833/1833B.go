package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution to problemB.txt for 1833B (Restore the Weather).
// We are given forecast temperatures array a and an array b with the actual
// temperatures in random order. We need to reorder b so that for every i the
// absolute difference |a_i - b_i| does not exceed k. It is guaranteed that such
// an ordering exists. The simple strategy is to sort the forecast values along
// with their original indices and also sort array b. Assign the j-th smallest b
// to the day with the j-th smallest a. This pairing always works under the
// given guarantee. Finally we output the rearranged b in the original order.
func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		_ = k // k is unused since sorting ensures a valid arrangement

		type pair struct{ val, idx int }
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &arr[i].val)
			arr[i].idx = i
		}

		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}

		sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
		sort.Ints(b)

		ans := make([]int, n)
		for i := 0; i < n; i++ {
			ans[arr[i].idx] = b[i]
		}

		for i := 0; i < n; i++ {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, ans[i])
		}
		fmt.Fprintln(writer)
	}
}
