package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, s int
		fmt.Fscan(in, &n, &s)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		used := make([]bool, n+1)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
			if b[i] != -1 {
				used[b[i]] = true
			}
		}
		ok := true
		for i := 0; i < n; i++ {
			if b[i] != -1 {
				if i+1 > b[i]+s {
					ok = false
					break
				}
			}
		}
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}
		avail := make([]int, 0, n)
		for v := 1; v <= n; v++ {
			if !used[v] {
				avail = append(avail, v)
			}
		}
		sort.Ints(avail)
		ans := int64(1)
		for i := n - 1; i >= 0; i-- {
			if b[i] == -1 {
				need := i + 1 - s
				if need < 1 {
					need = 1
				}
				idx := sort.SearchInts(avail, need)
				cnt := len(avail) - idx
				if cnt <= 0 {
					ans = 0
					break
				}
				ans = ans * int64(cnt) % mod
				avail = avail[:len(avail)-1]
			}
		}
		fmt.Fprintln(out, ans%mod)
	}
}
