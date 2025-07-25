package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
		var n int
		fmt.Fscan(reader, &n)
		size := 2 * n
		for i := 0; i < size; i++ {
			var sb strings.Builder
			sb.Grow(size)
			for j := 0; j < size; j++ {
				if ((i/2)+(j/2))%2 == 0 {
					sb.WriteByte('#')
				} else {
					sb.WriteByte('.')
				}
			}
			fmt.Fprintln(writer, sb.String())
		}
	}
}
