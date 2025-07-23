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

	var n, k int
	if _, err := fmt.Fscan(in, &n, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	sort.Ints(a)
	res := 0
	for i := 0; i < n; {
		val := a[i]
		j := i
		for j < n && a[j] == val {
			j++
		}
		if j < n && a[j] <= val+k {
			// all a[i:j] can be swallowed by some larger bacteria
		} else {
			res += j - i
		}
		i = j
	}
	fmt.Fprintln(out, res)
}
