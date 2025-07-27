package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var x, y int64
		fmt.Fscan(reader, &x, &y)
		limit := int64(math.Min(float64(y), math.Sqrt(float64(x))))
		var ans int64
		for k := int64(1); k <= limit; k++ {
			bmax := x/k - 1
			if bmax > y {
				bmax = y
			}
			if bmax > k {
				ans += bmax - k
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
