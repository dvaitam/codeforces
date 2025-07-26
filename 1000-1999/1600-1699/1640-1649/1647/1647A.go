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
		m := n / 3
		r := n % 3
		var sb strings.Builder
		switch r {
		case 1:
			sb.WriteByte('1')
		case 2:
			sb.WriteByte('2')
		}
		for i := 0; i < m; i++ {
			if r == 2 {
				sb.WriteString("12")
			} else {
				sb.WriteString("21")
			}
		}
		fmt.Fprintln(writer, sb.String())
	}
}
