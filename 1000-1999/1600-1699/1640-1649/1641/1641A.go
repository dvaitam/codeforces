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
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var x int64
		fmt.Fscan(in, &n, &x)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
		freq := make(map[int64]int)
		for _, v := range arr {
			freq[v]++
		}
		ans := 0
		for _, v := range arr {
			if freq[v] == 0 {
				continue
			}
			freq[v]--
			target := v * x
			if freq[target] > 0 {
				freq[target]--
			} else {
				ans++
			}
		}
		fmt.Fprintln(out, ans)
	}
}
