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
		var a, b string
		fmt.Fscan(reader, &a, &b)
		la := len(a)
		lb := len(b)
		ca := a[la-1]
		cb := b[lb-1]
		rank := func(c byte) int {
			if c == 'S' {
				return 0
			}
			if c == 'M' {
				return 1
			}
			return 2 // 'L'
		}
		ra := rank(ca)
		rb := rank(cb)
		if ra != rb {
			if ra > rb {
				fmt.Fprintln(writer, ">")
			} else {
				fmt.Fprintln(writer, "<")
			}
			continue
		}
		if ca == 'M' {
			fmt.Fprintln(writer, "=")
			continue
		}
		if ca == 'S' {
			if la == lb {
				fmt.Fprintln(writer, "=")
			} else if la < lb {
				fmt.Fprintln(writer, ">")
			} else {
				fmt.Fprintln(writer, "<")
			}
		} else { // 'L'
			if la == lb {
				fmt.Fprintln(writer, "=")
			} else if la > lb {
				fmt.Fprintln(writer, ">")
			} else {
				fmt.Fprintln(writer, "<")
			}
		}
	}
}
