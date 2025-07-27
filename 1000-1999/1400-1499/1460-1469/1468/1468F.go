package main

import (
	"bufio"
	"fmt"
	"os"
)

func gcd(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		type dir struct{ x, y int64 }
		cnt := make(map[dir]int64)
		for i := 0; i < n; i++ {
			var x, y, u, v int64
			fmt.Fscan(reader, &x, &y, &u, &v)
			dx := u - x
			dy := v - y
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			cnt[dir{dx, dy}]++
		}
		var ans int64
		for d, c := range cnt {
			opp := dir{-d.x, -d.y}
			if v, ok := cnt[opp]; ok {
				ans += c * v
			}
		}
		ans /= 2
		fmt.Fprintln(writer, ans)
	}
}
