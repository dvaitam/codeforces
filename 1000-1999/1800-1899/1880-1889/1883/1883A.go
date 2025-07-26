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
		pos := 1
		ans := 0
		for _, ch := range s {
			d := int(ch - '0')
			if d > pos {
				ans += d - pos
			} else {
				ans += pos - d
			}
			ans++
			pos = d
		}
		fmt.Fprintln(writer, ans)
	}
}
