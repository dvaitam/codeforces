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

	var a, b int64
	if _, err := fmt.Fscan(reader, &a, &b); err != nil {
		return
	}
	for {
		if a == 0 || b == 0 {
			break
		}
		if a >= 2*b {
			a %= 2 * b
			continue
		}
		if b >= 2*a {
			b %= 2 * a
			continue
		}
		break
	}
	fmt.Fprintf(writer, "%d %d\n", a, b)
}
