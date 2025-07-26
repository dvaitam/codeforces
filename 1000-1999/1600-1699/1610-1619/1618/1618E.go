package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	br := bufio.NewReader(os.Stdin)
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()

	var t int
	fmt.Fscan(br, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(br, &n)
		b := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			fmt.Fscan(br, &b[i])
			sum += b[i]
		}
		total := int64(n*(n+1)) / 2
		if sum*2%total != 0 {
			fmt.Fprintln(bw, "NO")
			continue
		}
		S := sum * 2 / total
		ans := make([]int64, n)
		ok := true
		for i := 0; i < n; i++ {
			prev := b[(i-1+n)%n]
			diff := b[i] - prev
			val := S - diff
			if val <= 0 || val%int64(n) != 0 {
				ok = false
				break
			}
			ans[i] = val / int64(n)
		}
		if !ok {
			fmt.Fprintln(bw, "NO")
			continue
		}
		fmt.Fprintln(bw, "YES")
		for i, v := range ans {
			if i > 0 {
				bw.WriteByte(' ')
			}
			fmt.Fprint(bw, v)
		}
		bw.WriteByte('\n')
	}
}
