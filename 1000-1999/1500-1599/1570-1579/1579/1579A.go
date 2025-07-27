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
		countA, countB, countC := 0, 0, 0
		for _, ch := range s {
			switch ch {
			case 'A':
				countA++
			case 'B':
				countB++
			case 'C':
				countC++
			}
		}
		if countB == countA+countC {
			writer.WriteString("YES\n")
		} else {
			writer.WriteString("NO\n")
		}
	}
}
