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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(reader, &s)
		if s[0] == '0' {
			fmt.Fprintln(writer, 0)
			continue
		}
		ans := 1
		if s[0] == '?' {
			ans *= 9
		}
		for i := 1; i < len(s); i++ {
			if s[i] == '?' {
				ans *= 10
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
