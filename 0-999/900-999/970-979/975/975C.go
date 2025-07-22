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

	var n, q int
	fmt.Fscan(in, &n, &q)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	pref := make([]int64, n)
	var sum int64
	for i := 0; i < n; i++ {
		sum += a[i]
		pref[i] = sum
	}
	var dmg int64
	for i := 0; i < q; i++ {
		var k int64
		fmt.Fscan(in, &k)
		dmg += k
		if dmg >= pref[n-1] {
			dmg = 0
			fmt.Fprintln(out, n)
			continue
		}
		idx := sort.Search(len(pref), func(i int) bool { return pref[i] > dmg })
		fmt.Fprintln(out, n-idx)
	}
}
