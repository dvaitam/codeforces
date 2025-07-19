package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	rdr := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	var T int
	fmt.Fscan(rdr, &T)
	for T > 0 {
		T--
		var n int
		fmt.Fscan(rdr, &n)
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(rdr, &x)
			if i%2 == 1 {
				if x > 0 {
					x = -x
				}
			} else {
				if x < 0 {
					x = -x
				}
			}
			fmt.Fprint(w, x)
			if i != n-1 {
				w.WriteByte(' ')
			}
		}
		w.WriteByte('\n')
	}
}
