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
		var n, k, q int
		fmt.Fscan(in, &n, &k, &q)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(in, &a[i])
		}

		diff := make([]int, n+1)
		for i := 1; i <= n; i++ {
			diff[i] = a[i] - i
		}

		bestLen := n - k + 1
		ans := make([]int, bestLen+1)

		count := make(map[int]int)
		freqCnt := make([]int, k+2)
		maxFreq := 0

		add := func(val int) {
			old := count[val]
			if old > 0 {
				freqCnt[old]--
			}
			count[val] = old + 1
			freqCnt[old+1]++
			if old+1 > maxFreq {
				maxFreq = old + 1
			}
		}

		remove := func(val int) {
			old := count[val]
			freqCnt[old]--
			if old == 1 {
				delete(count, val)
			} else {
				count[val] = old - 1
				freqCnt[old-1]++
			}
			for maxFreq > 0 && freqCnt[maxFreq] == 0 {
				maxFreq--
			}
		}

		for i := 1; i <= k; i++ {
			add(diff[i])
		}
		ans[1] = k - maxFreq

		for start := 2; start <= bestLen; start++ {
			remove(diff[start-1])
			add(diff[start+k-1])
			ans[start] = k - maxFreq
		}

		for i := 0; i < q; i++ {
			var l, r int
			fmt.Fscan(in, &l, &r)
			fmt.Fprintln(out, ans[l])
		}
	}
}
