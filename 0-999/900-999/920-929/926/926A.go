package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var l, r int64
	if _, err := fmt.Fscan(reader, &l, &r); err != nil {
		return
	}

	count := 0
	for p2 := int64(1); p2 <= r; p2 *= 2 {
		for p3 := int64(1); p2*p3 <= r; p3 *= 3 {
			val := p2 * p3
			if val >= l {
				count++
			}
		}
	}

	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, count)
}
