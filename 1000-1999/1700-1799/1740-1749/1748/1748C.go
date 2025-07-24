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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := range a {
			fmt.Fscan(in, &a[i])
		}

		pref := make([]int64, n+1)
		for i := 0; i < n; i++ {
			pref[i+1] = pref[i] + a[i]
		}

		zeros := make([]int, 0)
		for i, v := range a {
			if v == 0 {
				zeros = append(zeros, i+1) // 1-indexed
			}
		}

		if len(zeros) == 0 {
			cnt := 0
			for i := 1; i <= n; i++ {
				if pref[i] == 0 {
					cnt++
				}
			}
			fmt.Fprintln(out, cnt)
			continue
		}

		ans := 0
		firstZero := zeros[0]
		for i := 1; i < firstZero; i++ {
			if pref[i] == 0 {
				ans++
			}
		}

		pos := firstZero
		idx := 1
		for {
			nextZero := n + 1
			if idx < len(zeros) {
				nextZero = zeros[idx]
			}
			freq := make(map[int64]int)
			for j := pos + 1; j < nextZero; j++ {
				v := pref[j] - pref[pos]
				freq[v]++
			}
			best := freq[0] + 1
			for _, c := range freq {
				if c > best {
					best = c
				}
			}
			ans += best
			if nextZero == n+1 {
				break
			}
			pos = nextZero
			idx++
		}

		fmt.Fprintln(out, ans)
	}
}
