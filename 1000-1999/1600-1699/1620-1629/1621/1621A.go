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
		if k > (n+1)/2 {
			fmt.Fprintln(writer, -1)
			continue
		}
		for i := 0; i < n; i++ {
			line := make([]byte, n)
			for j := 0; j < n; j++ {
				line[j] = '.'
			}
			if i%2 == 0 && i/2 < k {
				line[i] = 'R'
			}
			fmt.Fprintln(writer, string(line))
		}
	}
}
