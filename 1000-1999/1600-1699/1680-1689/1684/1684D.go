package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Solution for problem described in problemD.txt in folder 1684.
// It selects up to k traps to jump so that the total damage is minimal.
// The strategy is to compute for each trap i a value a[i]-(n-i-1).
// Jumps are made at the k positions with the highest values.
// The minimal damage is then obtained using a closed form formula.

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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		type pair struct {
			val int64
			idx int
		}
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			arr[i] = pair{a[i] - int64(n-i-1), i}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].val > arr[j].val })

		var sumV int64
		for i := 0; i < k; i++ {
			sumV += arr[i].val
		}

		var total int64
		for i := 0; i < n; i++ {
			total += a[i]
		}

		ans := total - sumV - int64(k*(k-1)/2)
		fmt.Fprintln(writer, ans)
	}
}
