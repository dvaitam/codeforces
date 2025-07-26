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
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		hasOne := false
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] == 1 {
				hasOne = true
			}
		}
		if !hasOne {
			fmt.Fprintln(out, "YES")
			continue
		}
		sort.Ints(a)
		ok := true
		for i := 1; i < n; i++ {
			if a[i]-a[i-1] == 1 {
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
