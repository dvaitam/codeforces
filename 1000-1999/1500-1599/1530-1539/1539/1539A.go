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

	var k int
	if _, err := fmt.Fscan(in, &k); err != nil {
		return
	}
	for ; k > 0; k-- {
		var n, x, t int64
		fmt.Fscan(in, &n, &x, &t)
		q := t / x
		if q > n-1 {
			q = n - 1
		}
		// total dissatisfaction = q*n - q*(q+1)/2
		ans := q*n - q*(q+1)/2
		fmt.Fprintln(out, ans)
	}
}
