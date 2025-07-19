package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t, n int
	fmt.Fscan(reader, &t)
	for t > 0 {
		t--
		fmt.Fscan(reader, &n)
		prev := int64(1)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(reader, &x)
			if x == 1 {
				x++
			}
			if x%prev == 0 {
				x++
			}
			if i > 0 {
				writer.WriteByte(' ')
			}
			writer.WriteString(strconv.FormatInt(x, 10))
			prev = x
		}
		writer.WriteByte('\n')
	}
}
