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

	var q int
	if _, err := fmt.Fscan(in, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(in, &n)
		m := 4 * n
		arr := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &arr[i])
		}
		sort.Ints(arr)
		area := arr[0] * arr[m-1]
		ok := true
		for i := 0; i < n; i++ {
			if arr[2*i] != arr[2*i+1] || arr[m-2*i-2] != arr[m-2*i-1] || arr[2*i]*arr[m-2*i-1] != area {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
