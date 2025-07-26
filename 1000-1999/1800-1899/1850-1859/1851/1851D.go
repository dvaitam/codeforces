package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		pref := make([]int64, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(reader, &pref[i])
		}
		total := int64(n) * int64(n+1) / 2
		if pref[n-2] > total {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if pref[n-2] < total {
			diffs := make([]int, 0, n)
			prev := int64(0)
			ok := true
			seen := make([]bool, n+1)
			for _, v := range pref {
				diff := int(v - prev)
				if diff < 1 || diff > n || seen[diff] {
					ok = false
					break
				}
				seen[diff] = true
				diffs = append(diffs, diff)
				prev = v
			}
			missing := int(total - pref[n-2])
			if missing < 1 || missing > n || seen[missing] {
				ok = false
			}
			if ok {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		} else {
			diffs := make([]int, 0, n-1)
			prev := int64(0)
			seen := make([]bool, n+1)
			large := -1
			for _, v := range pref {
				diff := int(v - prev)
				if diff >= 1 && diff <= n && !seen[diff] {
					seen[diff] = true
				} else {
					if large == -1 {
						large = diff
					} else {
						large = -2
					}
				}
				diffs = append(diffs, diff)
				prev = v
			}
			missing := []int{}
			for i := 1; i <= n; i++ {
				if !seen[i] {
					missing = append(missing, i)
				}
			}
			if large == -1 && len(missing) == 1 {
				// large difference might be implicitly the missing value
				large = missing[0]
				missing = missing[:0]
			}
			if len(missing) == 2 && large == missing[0]+missing[1] {
				fmt.Fprintln(writer, "YES")
			} else {
				fmt.Fprintln(writer, "NO")
			}
		}
	}
}
