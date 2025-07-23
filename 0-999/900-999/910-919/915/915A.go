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

	var n, k int
	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}

	best := int(1<<31 - 1) // large number
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(reader, &a)
		if k%a == 0 {
			hours := k / a
			if hours < best {
				best = hours
			}
		}
	}

	fmt.Fprintln(writer, best)
}
