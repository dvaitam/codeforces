package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int, 0, n)
		b := make([]int, 0, n)
		for i := 0; i < 2*n; i++ {
			var x, y int
			fmt.Fscan(reader, &x, &y)
			if x == 0 {
				if y < 0 {
					y = -y
				}
				a = append(a, y)
			} else {
				if x < 0 {
					x = -x
				}
				b = append(b, x)
			}
		}
		sort.Ints(a)
		sort.Ints(b)
		ans := 0.0
		for i := 0; i < n; i++ {
			ans += math.Hypot(float64(a[i]), float64(b[i]))
		}
		fmt.Printf("%.9f\n", ans)
	}
}
