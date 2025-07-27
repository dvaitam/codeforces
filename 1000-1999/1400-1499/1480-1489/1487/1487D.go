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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var n int64
		fmt.Fscan(reader, &n)
		m := int64(math.Sqrt(float64(2*n - 1)))
		ans := (m - 1) / 2
		if ans < 0 {
			ans = 0
		}
		fmt.Fprintln(writer, ans)
	}
}
