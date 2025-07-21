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
		p := make([]int, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &p[i])
		}
		cnt := 0
		for i := 1; i <= n; i++ {
			for j := i + 1; j <= n; j++ {
				if (p[i]*p[j])%(i*j) == 0 {
					cnt++
				}
			}
		}
		fmt.Fprintln(writer, cnt)
	}
}
