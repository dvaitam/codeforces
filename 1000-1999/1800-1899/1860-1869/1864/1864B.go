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

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		if k%2 == 0 {
			b := []byte(s)
			sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
			fmt.Fprintln(out, string(b))
		} else {
			odd := make([]byte, 0, (n+1)/2)
			even := make([]byte, 0, n/2)
			for i := 0; i < n; i++ {
				if i%2 == 0 {
					odd = append(odd, s[i])
				} else {
					even = append(even, s[i])
				}
			}
			sort.Slice(odd, func(i, j int) bool { return odd[i] < odd[j] })
			sort.Slice(even, func(i, j int) bool { return even[i] < even[j] })
			res := make([]byte, n)
			oi, ei := 0, 0
			for i := 0; i < n; i++ {
				if i%2 == 0 {
					res[i] = odd[oi]
					oi++
				} else {
					res[i] = even[ei]
					ei++
				}
			}
			fmt.Fprintln(out, string(res))
		}
	}
}
