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
		var n, m int
		var k, H int
		fmt.Fscan(reader, &n, &m, &k, &H)
		count := 0
		for i := 0; i < n; i++ {
			var h int
			fmt.Fscan(reader, &h)
			diff := int(math.Abs(float64(h - H)))
			if diff%k == 0 {
				d := diff / k
				if d >= 1 && d <= m-1 {
					count++
				}
			}
		}
		fmt.Fprintln(writer, count)
	}
}
