package main

import (
	"bufio"
	"fmt"
	"os"
)

func pow10(x int) int64 {
	res := int64(1)
	for i := 0; i < x; i++ {
		res *= 10
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		var k int64
		fmt.Fscan(reader, &n, &k)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}
		k++ // need at least k+1 banknotes
		var ans int64
		for i := 0; i < n; i++ {
			var limit int64
			if i+1 < n {
				diff := a[i+1] - a[i]
				limit = pow10(diff) - 1
			} else {
				limit = k
			}
			use := k
			if use > limit {
				use = limit
			}
			ans += use * pow10(a[i])
			k -= use
		}
		fmt.Fprintln(writer, ans)
	}
}
