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

	pow3 := make([]int64, 0)
	cost := make([]int64, 0)
	var p int64 = 1
	for i := 0; i <= 40; i++ {
		pow3 = append(pow3, p)
		var c int64
		if i == 0 {
			c = 3
		} else {
			c = pow3[i-1]*9 + int64(i)*pow3[i-1]
		}
		cost = append(cost, c)
		p *= 3
	}

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(in, &n)
		ans := int64(0)
		idx := 0
		for n > 0 {
			d := n % 3
			if d > 0 {
				ans += int64(d) * cost[idx]
			}
			n /= 3
			idx++
		}
		fmt.Fprintln(out, ans)
	}
}
