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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i], &b[i])
		}
		// compute additional bullets needed for each monster
		extra := make([]int64, n)
		var sum int64
		for i := 0; i < n; i++ {
			prev := (i + n - 1) % n
			if a[i] > b[prev] {
				extra[i] = a[i] - b[prev]
			} else {
				extra[i] = 0
			}
			sum += extra[i]
		}
		// choose starting monster that minimizes total bullets
		ans := sum + a[0] - extra[0]
		for i := 1; i < n; i++ {
			cand := sum + a[i] - extra[i]
			if cand < ans {
				ans = cand
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
