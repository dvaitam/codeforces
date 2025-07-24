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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Ints(a)
		var ans int64
		pref := 0
		for i := 0; i < n; {
			j := i
			for j < n && a[j] == a[i] {
				j++
			}
			pref += j - i
			if pref < n {
				v := int64(pref) * int64(n-pref)
				if v > ans {
					ans = v
				}
			}
			i = j
		}
		fmt.Fprintln(out, ans)
	}
}
