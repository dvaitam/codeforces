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

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var n int
		fmt.Fscan(reader, &n)
		var t string
		fmt.Fscan(reader, &t)
		res := make([]byte, 0, n)
		for i := n - 1; i >= 0; {
			if t[i] == '0' {
				num := int(t[i-2]-'0')*10 + int(t[i-1]-'0')
				res = append(res, byte('a'+num-1))
				i -= 3
			} else {
				num := int(t[i] - '0')
				res = append(res, byte('a'+num-1))
				i--
			}
		}
		// reverse result
		for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
			res[i], res[j] = res[j], res[i]
		}
		fmt.Fprintln(writer, string(res))
	}
}
