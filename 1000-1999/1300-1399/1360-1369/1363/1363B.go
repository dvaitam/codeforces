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
		n := len(s)
		prefix0 := make([]int, n+1)
		prefix1 := make([]int, n+1)
		for i := 0; i < n; i++ {
			prefix0[i+1] = prefix0[i]
			prefix1[i+1] = prefix1[i]
			if s[i] == '0' {
				prefix0[i+1]++
			} else {
				prefix1[i+1]++
			}
		}
		total0 := prefix0[n]
		total1 := prefix1[n]
		minOps := n
		for i := 0; i <= n; i++ {
			ops01 := prefix1[i] + (total0 - prefix0[i])
			if ops01 < minOps {
				minOps = ops01
			}
			ops10 := prefix0[i] + (total1 - prefix1[i])
			if ops10 < minOps {
				minOps = ops10
			}
		}
		fmt.Fprintln(writer, minOps)
	}
}
