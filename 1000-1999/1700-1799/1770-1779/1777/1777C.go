package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MaxA = 100000

var preDivs = make([][]int, MaxA+1)

func init() {
	for d := 1; d <= MaxA; d++ {
		for m := d; m <= MaxA; m += d {
			preDivs[m] = append(preDivs[m], d)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		sort.Ints(a)

		allCnt := make([]int, m+1)
		allCov := 0
		for i := 0; i < n; i++ {
			for _, d := range preDivs[a[i]] {
				if d > m {
					break
				}
				if allCnt[d] == 0 {
					allCov++
				}
				allCnt[d]++
			}
		}
		if allCov < m {
			fmt.Fprintln(writer, -1)
			continue
		}

		cnt := make([]int, m+1)
		covered := 0
		ans := 1<<31 - 1
		l := 0
		for r := 0; r < n; r++ {
			for _, d := range preDivs[a[r]] {
				if d > m {
					break
				}
				cnt[d]++
				if cnt[d] == 1 {
					covered++
				}
			}
			for covered == m {
				if diff := a[r] - a[l]; diff < ans {
					ans = diff
				}
				for _, d := range preDivs[a[l]] {
					if d > m {
						break
					}
					cnt[d]--
					if cnt[d] == 0 {
						covered--
					}
				}
				l++
			}
		}
		if ans == 1<<31-1 {
			fmt.Fprintln(writer, -1)
		} else {
			fmt.Fprintln(writer, ans)
		}
	}
}
