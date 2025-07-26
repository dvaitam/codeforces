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

		hasCaret := false
		for i := 0; i < len(s); i++ {
			if s[i] == '^' {
				hasCaret = true
				break
			}
		}

		if !hasCaret {
			fmt.Fprintln(writer, len(s)+1)
			continue
		}
		if len(s) == 1 {
			fmt.Fprintln(writer, 1)
			continue
		}
		ans := 0
		if s[0] == '_' {
			ans++
		}
		if s[len(s)-1] == '_' {
			ans++
		}
		for i := 1; i < len(s); i++ {
			if s[i] == '_' && s[i-1] == '_' {
				ans++
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
