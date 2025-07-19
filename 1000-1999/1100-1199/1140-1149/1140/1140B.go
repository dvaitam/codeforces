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
	for i := 0; i < t; i++ {
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		res1 := 0
		for j := 0; j < n; j++ {
			if s[j] != '<' {
				break
			}
			res1++
		}
		res2 := 0
		for j := n - 1; j >= 0; j-- {
			if s[j] != '>' {
				break
			}
			res2++
		}
		if res1 < res2 {
			fmt.Fprintln(writer, res1)
		} else {
			fmt.Fprintln(writer, res2)
		}
	}
}
