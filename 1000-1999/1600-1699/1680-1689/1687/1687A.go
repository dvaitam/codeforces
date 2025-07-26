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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(reader, &n, &k)
		a := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
			sum += a[i]
		}
		if k < n {
			pref := make([]int64, n+1)
			for i := 0; i < n; i++ {
				pref[i+1] = pref[i] + a[i]
			}
			var best int64
			for i := 0; i+k <= n; i++ {
				s := pref[i+k] - pref[i]
				if i == 0 || s > best {
					best = s
				}
			}
			k64 := int64(k)
			fmt.Fprintln(writer, best+k64*(k64-1)/2)
		} else {
			k64 := int64(k)
			n64 := int64(n)
			add := n64*(n64-1)/2 + (k64-n64)*n64
			fmt.Fprintln(writer, sum+add)
		}
	}
}
