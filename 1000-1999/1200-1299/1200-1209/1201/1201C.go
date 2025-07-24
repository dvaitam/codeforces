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

	var n int
	var k int64
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}

	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &arr[i])
	}

	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	mid := n / 2

	lo := arr[mid]
	hi := arr[mid] + k + 1
	for lo+1 < hi {
		m := (lo + hi) / 2
		var need int64
		for i := mid; i < n; i++ {
			if arr[i] < m {
				need += m - arr[i]
				if need > k {
					break
				}
			}
		}
		if need <= k {
			lo = m
		} else {
			hi = m
		}
	}

	fmt.Fprintln(out, lo)
}
