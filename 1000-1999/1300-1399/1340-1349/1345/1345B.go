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

	cards := make([]int, 0)
	for h := 1; ; h++ {
		val := (3*h*h + h) / 2
		if val > 1e9 {
			break
		}
		cards = append(cards, val)
	}

	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		count := 0
		for n >= 2 {
			l, r := 0, len(cards)-1
			best := -1
			for l <= r {
				mid := (l + r) / 2
				if cards[mid] <= n {
					best = cards[mid]
					l = mid + 1
				} else {
					r = mid - 1
				}
			}
			if best == -1 {
				break
			}
			n -= best
			count++
		}
		fmt.Fprintln(writer, count)
	}
}
