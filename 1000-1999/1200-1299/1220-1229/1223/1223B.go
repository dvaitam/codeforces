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
	if _, err := fmt.Fscan(reader, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		letters := [26]bool{}
		for _, ch := range s {
			letters[ch-'a'] = true
		}
		ok := false
		for _, ch := range t {
			if letters[ch-'a'] {
				ok = true
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
