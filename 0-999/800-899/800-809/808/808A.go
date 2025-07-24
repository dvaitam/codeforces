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

	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var best int64 = 1<<63 - 1
	for d := int64(1); d <= 9; d++ {
		x := d
		for k := 0; k <= 10; k++ {
			if x > n {
				if x-n < best {
					best = x - n
				}
				break
			}
			x *= 10
		}
	}
	fmt.Fprintln(writer, best)
}
