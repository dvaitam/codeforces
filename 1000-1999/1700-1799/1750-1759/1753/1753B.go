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

	var n, x int
	if _, err := fmt.Fscan(in, &n, &x); err != nil {
		return
	}
	cnt := make([]int64, x+1)
	for i := 0; i < n; i++ {
		var a int
		fmt.Fscan(in, &a)
		if a <= x {
			cnt[a]++
		}
	}
	q := cnt[1]
	for i := 2; i <= x; i++ {
		val := q + cnt[i]*int64(i)
		if val%int64(i) != 0 {
			fmt.Fprintln(out, "No")
			return
		}
		q = val / int64(i)
	}
	fmt.Fprintln(out, "Yes")
}
