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
		var n, q int
		fmt.Fscan(in, &n, &q)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		pref := make([]int, n+1)
		for i := 1; i <= n; i++ {
			pref[i] = pref[i-1] ^ a[i-1]
		}
		pos := make(map[int][]int, n+1)
		for i := 0; i <= n; i++ {
			v := pref[i]
			pos[v] = append(pos[v], i)
		}
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			total := pref[r] ^ pref[l-1]
			if total == 0 {
				fmt.Fprintln(out, "YES")
				continue
			}
			list1 := pos[pref[r]]
			i := sort.Search(len(list1), func(i int) bool { return list1[i] >= l })
			if i == len(list1) || list1[i] >= r {
				fmt.Fprintln(out, "NO")
				continue
			}
			j1 := list1[i]
			list2 := pos[pref[l-1]]
			j := sort.Search(len(list2), func(i int) bool { return list2[i] > j1 })
			if j == len(list2) || list2[j] >= r {
				fmt.Fprintln(out, "NO")
			} else {
				fmt.Fprintln(out, "YES")
			}
		}
	}
}
