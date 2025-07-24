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
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &b[i])
		}
		ans := 0
		for bit := 29; bit >= 0; bit-- {
			cand := ans | (1 << uint(bit))
			count := make(map[int]int)
			for i := 0; i < n; i++ {
				count[a[i]&cand]++
			}
			ok := true
			for i := 0; i < n; i++ {
				key := cand ^ (b[i] & cand)
				if count[key] == 0 {
					ok = false
					break
				}
				count[key]--
			}
			if ok {
				ans = cand
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
