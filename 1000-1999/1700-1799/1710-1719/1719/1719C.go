package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, q int
		fmt.Fscan(reader, &n, &q)
		a := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		// index of global maximum (strength n)
		maxIdx := 1
		for i := 1; i <= n; i++ {
			if a[i] > a[maxIdx] {
				maxIdx = i
			}
		}

		wins := make([][]int, n+1)
		champion := 1
		for j := 2; j <= n; j++ {
			round := j - 1
			if a[champion] > a[j] {
				wins[champion] = append(wins[champion], round)
			} else {
				wins[j] = append(wins[j], round)
				champion = j
			}
		}

		for ; q > 0; q-- {
			var idx int
			var k int
			fmt.Fscan(reader, &idx, &k)
			// number of wins in rounds <= k (but rounds start from 1)
			limit := k
			if limit > n-1 {
				limit = n - 1
			}
			// count wins up to 'limit'
			rounds := wins[idx]
			cnt := sort.Search(len(rounds), func(i int) bool { return rounds[i] > limit })
			if idx == maxIdx && k > n-1 {
				cnt += k - (n - 1)
			}
			fmt.Fprintln(writer, cnt)
		}
	}
}
