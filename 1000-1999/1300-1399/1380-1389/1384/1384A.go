package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		ar := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &ar[i])
		}
		sam := make([]byte, 101)
		for i := range sam {
			sam[i] = 'a'
		}
		// initial string
		out.Write(sam)
		out.WriteByte('\n')
		for i := 0; i < n; i++ {
			idx := ar[i]
			if sam[idx] == 'a' {
				sam[idx] = 'b'
			} else {
				sam[idx] = 'a'
			}
			out.Write(sam)
			out.WriteByte('\n')
		}
	}
}
