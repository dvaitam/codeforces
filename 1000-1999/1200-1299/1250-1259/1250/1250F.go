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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	best := int(1<<31 - 1)
	for i := 1; i*i <= n; i++ {
		if n%i == 0 {
			j := n / i
			p := 2 * (i + j)
			if p < best {
				best = p
			}
		}
	}
	fmt.Fprintln(writer, best)
}
