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
		var n int
		fmt.Fscan(reader, &n)
		var s string
		fmt.Fscan(reader, &s)

		possible := true
		i := 0
		for i < n {
			if s[i] == 'W' {
				i++
				continue
			}
			j := i
			hasR := false
			hasB := false
			for j < n && s[j] != 'W' {
				if s[j] == 'R' {
					hasR = true
				} else if s[j] == 'B' {
					hasB = true
				}
				j++
			}
			if !(hasR && hasB) {
				possible = false
				break
			}
			i = j
		}
		if possible {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
