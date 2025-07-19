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
	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(reader, &n)
		ev := make([]int64, 0, n)
		for i := 0; i < n; i++ {
			var a int64
			fmt.Fscan(reader, &a)
			if a%2 != 0 {
				fmt.Fprint(writer, a, " ")
			} else {
				ev = append(ev, a)
			}
		}
		for _, v := range ev {
			fmt.Fprint(writer, v, " ")
		}
		fmt.Fprintln(writer)
	}
}
