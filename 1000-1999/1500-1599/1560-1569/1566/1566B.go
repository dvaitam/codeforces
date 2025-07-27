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
		var s string
		fmt.Fscan(reader, &s)
		zeroGroups := 0
		prev := byte('1')
		for i := 0; i < len(s); i++ {
			if s[i] == '0' && prev != '0' {
				zeroGroups++
			}
			prev = s[i]
		}
		if zeroGroups > 2 {
			zeroGroups = 2
		}
		fmt.Fprintln(writer, zeroGroups)
	}
}
