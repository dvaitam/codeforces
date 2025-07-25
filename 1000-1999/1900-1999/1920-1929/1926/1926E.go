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
		var n, k int64
		fmt.Fscan(in, &n, &k)
		pow := int64(1)
		for {
			count := n/pow - n/(pow*2)
			if k > count {
				k -= count
				pow <<= 1
			} else {
				ans := pow * (2*k - 1)
				fmt.Fprintln(out, ans)
				break
			}
		}
	}
}
