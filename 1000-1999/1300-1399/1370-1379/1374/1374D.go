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
		var k int64
		fmt.Fscan(reader, &n, &k)
		freq := make(map[int64]int)
		for i := 0; i < n; i++ {
			var a int64
			fmt.Fscan(reader, &a)
			r := a % k
			if r != 0 {
				d := k - r
				freq[d]++
			}
		}
		var ans int64
		for d, c := range freq {
			cand := int64(c-1)*k + int64(d)
			if cand > ans {
				ans = cand
			}
		}
		if ans == 0 {
			fmt.Fprintln(writer, 0)
		} else {
			fmt.Fprintln(writer, ans+1)
		}
	}
}
