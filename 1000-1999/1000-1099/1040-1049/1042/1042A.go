package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int64
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	var maxx, sum int64
	for i := int64(0); i < n; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		if x > maxx {
			maxx = x
		}
		sum += x
	}
	// Determine minimum possible maximum after distributing m
	var u int64
	if (sum+m)%n != 0 {
		u = 1
	}
	// Average (rounded up) vs current max
	res1 := (sum+m)/n + u
	if res1 < maxx {
		res1 = maxx
	}
	// Maximum possible maximum
	res2 := maxx + m
	fmt.Fprintln(writer, res1, res2)
}
