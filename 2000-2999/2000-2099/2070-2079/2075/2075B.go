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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		sort.Slice(arr, func(i, j int) bool {
			return arr[i] > arr[j]
		})

		var topSum int64
		for i := 0; i < k; i++ {
			topSum += arr[i]
		}

		var best int64
		for j := 0; j < n; j++ {
			var sumOthers int64
			if j < k {
				sumOthers = topSum - arr[j]
				sumOthers += arr[k]
			} else {
				sumOthers = topSum
			}
			candidate := sumOthers + arr[j]
			if candidate > best {
				best = candidate
			}
		}

		fmt.Fprintln(out, best)
	}
}
