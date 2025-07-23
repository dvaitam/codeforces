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

	var n int
	var d int64
	fmt.Fscan(reader, &n, &d)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}

	maxAfter := make([]int64, n+1)
	maxAfter[n] = pref[n]
	for i := n - 1; i >= 0; i-- {
		if pref[i] > maxAfter[i+1] {
			maxAfter[i] = pref[i]
		} else {
			maxAfter[i] = maxAfter[i+1]
		}
	}

	shift := int64(0)
	count := 0
	for i := 1; i <= n; i++ {
		cur := pref[i] + shift
		if cur > d {
			fmt.Fprintln(writer, -1)
			return
		}
		if a[i-1] == 0 && cur < 0 {
			limit := d - maxAfter[i]
			if limit < -pref[i] {
				fmt.Fprintln(writer, -1)
				return
			}
			if shift < limit {
				shift = limit
				count++
			}
			cur = pref[i] + shift
			if cur < 0 || cur > d {
				fmt.Fprintln(writer, -1)
				return
			}
		}
	}
	fmt.Fprintln(writer, count)
}
