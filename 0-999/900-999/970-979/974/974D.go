package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod int64 = 1000000007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var res int64
		if n == 0 {
			res = 0
		} else if n == 1 {
			res = 1
		} else {
			a, b := int64(0), int64(1)
			for i := 2; i <= n; i++ {
				a, b = b, (a+b)%mod
			}
			res = b
		}
		fmt.Fprintln(out, res)
	}
}
