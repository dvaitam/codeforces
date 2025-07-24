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
		var a [2][2]int
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				fmt.Fscan(in, &a[i][j])
			}
		}
		count := a[0][0] + a[0][1] + a[1][0] + a[1][1]
		switch count {
		case 0:
			fmt.Fprintln(out, 0)
		case 4:
			fmt.Fprintln(out, 2)
		default:
			fmt.Fprintln(out, 1)
		}
	}
}
