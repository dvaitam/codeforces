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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
		if k >= 3 {
			fmt.Fprintln(out, 0)
			continue
		}
		minVal := a[0]
		hasDup := false
		diff1 := int64(1 << 62)
		for i := 1; i < n; i++ {
			if a[i] == a[i-1] {
				hasDup = true
			}
			d := a[i] - a[i-1]
			if d < diff1 {
				diff1 = d
			}
		}
		if k == 0 {
			fmt.Fprintln(out, minVal)
			continue
		}
		if k == 1 {
			if hasDup {
				fmt.Fprintln(out, 0)
			} else if minVal < diff1 {
				fmt.Fprintln(out, minVal)
			} else {
				fmt.Fprintln(out, diff1)
			}
			continue
		}
		// k == 2
		if hasDup {
			fmt.Fprintln(out, 0)
			continue
		}
		diff2 := int64(1 << 62)
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				d := a[j] - a[i]
				idx := sort.Search(len(a), func(x int) bool { return a[x] >= d })
				if idx < n {
					v := a[idx] - d
					if v < diff2 {
						diff2 = v
					}
					if diff2 == 0 {
						break
					}
				}
				if idx > 0 {
					v := d - a[idx-1]
					if v < diff2 {
						diff2 = v
					}
					if diff2 == 0 {
						break
					}
				}
			}
			if diff2 == 0 {
				break
			}
		}
		ans := minVal
		if diff1 < ans {
			ans = diff1
		}
		if diff2 < ans {
			ans = diff2
		}
		fmt.Fprintln(out, ans)
	}
}
