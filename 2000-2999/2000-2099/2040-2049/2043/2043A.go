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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		k := 0
		for n > 3 {
			k++
			n /= 4
		}
		res := int64(1)
		for i := 0; i < k; i++ {
			res <<= 1
		}
		fmt.Fprintln(out, res)
	}
}
