package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var r, d int
	if _, err := fmt.Fscan(in, &r, &d); err != nil {
		return
	}
	var n int
	fmt.Fscan(in, &n)
	outer := float64(r)
	inner := float64(r - d)
	ans := 0
	for i := 0; i < n; i++ {
		var x, y, rr int
		fmt.Fscan(in, &x, &y, &rr)
		dist := math.Hypot(float64(x), float64(y))
		fr := float64(rr)
		if dist+fr <= outer && dist-fr >= inner {
			ans++
		}
	}
	fmt.Println(ans)
}
