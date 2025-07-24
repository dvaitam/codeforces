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

	var n, k, q int
	if _, err := fmt.Fscan(in, &n, &k, &q); err != nil {
		return
	}

	const maxTemp = 200000
	diff := make([]int, maxTemp+2)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		diff[l]++
		if r+1 <= maxTemp {
			diff[r+1]--
		}
	}

	freq := make([]int, maxTemp+1)
	for i := 1; i <= maxTemp; i++ {
		freq[i] = freq[i-1] + diff[i]
	}

	pref := make([]int, maxTemp+1)
	for i := 1; i <= maxTemp; i++ {
		pref[i] = pref[i-1]
		if freq[i] >= k {
			pref[i]++
		}
	}

	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a < 1 {
			a = 1
		}
		if b > maxTemp {
			b = maxTemp
		}
		if a > b {
			fmt.Fprintln(out, 0)
			continue
		}
		fmt.Fprintln(out, pref[b]-pref[a-1])
	}
}
