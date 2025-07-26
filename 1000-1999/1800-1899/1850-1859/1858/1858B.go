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
		var n, m, d int
		fmt.Fscan(reader, &n, &m, &d)
		s := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &s[i])
		}

		base := 1
		prev := 1
		for _, x := range s {
			if x == 1 {
				prev = x
				continue
			}
			base += (x - prev - 1) / d
			base++
			prev = x
		}
		base += (n - prev) / d

		best := int(1<<63 - 1)
		count := 0
		for i, x := range s {
			var newTotal int
			if x == 1 {
				newTotal = base
			} else {
				prevPos := 1
				if i > 0 {
					prevPos = s[i-1]
				}
				nextPos := n + 1
				if i < m-1 {
					nextPos = s[i+1]
				}
				gap1 := (x - prevPos - 1) / d
				gap2 := (nextPos - x - 1) / d
				combined := (nextPos - prevPos - 1) / d
				newTotal = base - 1 - gap1 - gap2 + combined
			}
			if newTotal < best {
				best = newTotal
				count = 1
			} else if newTotal == best {
				count++
			}
		}
		fmt.Fprintf(writer, "%d %d\n", best, count)
	}
}
