package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		l := make([]int, n)
		r := make([]int, n)
		maxVal := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &l[i], &r[i])
			if r[i] > maxVal {
				maxVal = r[i]
			}
		}

		singletonCount := make([]int, maxVal+2)
		for i := 0; i < n; i++ {
			if l[i] == r[i] {
				singletonCount[l[i]]++
			}
		}

		blockedPref := make([]int, maxVal+2)
		for v := 1; v <= maxVal; v++ {
			blockedPref[v] = blockedPref[v-1]
			if singletonCount[v] > 0 {
				blockedPref[v]++
			}
		}

		var builder strings.Builder
		builder.Grow(n)
		for i := 0; i < n; i++ {
			if l[i] == r[i] {
				if singletonCount[l[i]] == 1 {
					builder.WriteByte('1')
				} else {
					builder.WriteByte('0')
				}
				continue
			}
			blocked := blockedPref[r[i]] - blockedPref[l[i]-1]
			length := r[i] - l[i] + 1
			if blocked < length {
				builder.WriteByte('1')
			} else {
				builder.WriteByte('0')
			}
		}
		fmt.Fprintln(out, builder.String())
	}
}
