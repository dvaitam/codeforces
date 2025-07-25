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
		var n, k, d int
		fmt.Fscan(reader, &n, &k, &d)
		shows := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &shows[i])
		}
		cnt := make([]int, k+1)
		unique := 0
		for i := 0; i < d; i++ {
			if cnt[shows[i]] == 0 {
				unique++
			}
			cnt[shows[i]]++
		}
		best := unique
		for i := d; i < n; i++ {
			if cnt[shows[i]] == 0 {
				unique++
			}
			cnt[shows[i]]++
			rem := shows[i-d]
			cnt[rem]--
			if cnt[rem] == 0 {
				unique--
			}
			if unique < best {
				best = unique
			}
		}
		fmt.Fprintln(writer, best)
	}
}
