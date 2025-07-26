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
		var s, c string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &c)
		res := "NO"
		for i := 0; i < len(s); i++ {
			if s[i] == c[0] && i%2 == 0 {
				res = "YES"
				break
			}
		}
		fmt.Fprintln(writer, res)
	}
}
