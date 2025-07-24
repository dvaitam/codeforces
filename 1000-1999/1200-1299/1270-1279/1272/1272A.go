package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		best := int(^uint(0) >> 1) // max int
		for da := -1; da <= 1; da++ {
			for db := -1; db <= 1; db++ {
				for dc := -1; dc <= 1; dc++ {
					aa := a + da
					bb := b + db
					cc := c + dc
					dist := abs(aa-bb) + abs(aa-cc) + abs(bb-cc)
					if dist < best {
						best = dist
					}
				}
			}
		}
		fmt.Fprintln(writer, best)
	}
}
