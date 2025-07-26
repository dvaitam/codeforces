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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		mism := 0
		for i := 0; i < n; i++ {
			if p[i]%k != (i+1)%k {
				mism++
			}
		}
		if mism == 0 {
			fmt.Fprintln(writer, 0)
		} else if mism == 2 {
			fmt.Fprintln(writer, 1)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
