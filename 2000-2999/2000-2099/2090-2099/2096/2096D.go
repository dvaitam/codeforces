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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)

		colParity := make(map[int64]byte)
		diagParity := make(map[int64]byte) // key is x+y

		for i := 0; i < n; i++ {
			var x, y int64
			fmt.Fscan(in, &x, &y)
			colParity[x] ^= 1
			diagParity[x+y] ^= 1
		}

		var s, diag int64
		for k, v := range colParity {
			if v&1 == 1 {
				s = k
				break
			}
		}
		for k, v := range diagParity {
			if v&1 == 1 {
				diag = k
				break
			}
		}
		t := diag - s
		fmt.Fprintln(out, s, t)
	}
}
