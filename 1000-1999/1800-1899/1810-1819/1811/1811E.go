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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var k int64
		fmt.Fscan(reader, &k)
		if k == 0 {
			fmt.Fprintln(writer, 0)
			continue
		}
		var digits []byte
		for k > 0 {
			d := k % 9
			if d >= 4 {
				digits = append(digits, byte('0'+d+1))
			} else {
				digits = append(digits, byte('0'+d))
			}
			k /= 9
		}
		// reverse digits
		for i, j := 0, len(digits)-1; i < j; i, j = i+1, j-1 {
			digits[i], digits[j] = digits[j], digits[i]
		}
		writer.Write(digits)
		writer.WriteByte('\n')
	}
}
