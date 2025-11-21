package main

import (
	"bufio"
	"fmt"
	"os"
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
		var n, k int
		fmt.Fscan(in, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		freq := make([]int, n+1)
		for _, v := range a {
			freq[v]++
		}

		ok := true
		for _, f := range freq {
			if f%k != 0 {
				ok = false
				break
			}
		}
		if !ok {
			fmt.Fprintln(out, 0)
			continue
		}

		limit := make([]int, n+1)
		for val, f := range freq {
			limit[val] = f / k
		}

		cnt := make([]int, n+1)
		left := 0
		var ans int64
		for right := 0; right < n; right++ {
			v := a[right]
			cnt[v]++
			for cnt[v] > limit[v] {
				cnt[a[left]]--
				left++
			}
			ans += int64(right - left + 1)
		}
		fmt.Fprintln(out, ans)
	}
}
