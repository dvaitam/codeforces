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

	var T int
	if _, err := fmt.Fscan(reader, &T); err != nil {
		return
	}
	for T > 0 {
		T--
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n+1)
		var sum int64
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &a[i])
			sum += a[i]
		}
		if sum%int64(n) != 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		avg := sum / int64(n)
		// total operations: 3*(n-1)
		fmt.Fprintln(writer, 3*(n-1))
		// make each a[i] divisible by i and move to a[1]
		for i := 2; i <= n; i++ {
			rem := a[i] % int64(i)
			delta := (int64(i) - rem) % int64(i)
			// move delta from 1 to i
			fmt.Fprintf(writer, "1 %d %d\n", i, delta)
			a[1] -= delta
			a[i] += delta
			// move a[i] from i to 1 in units of size i
			x := a[i] / int64(i)
			fmt.Fprintf(writer, "%d 1 %d\n", i, x)
			a[1] += a[i]
			a[i] = 0
		}
		// distribute avg to each a[i]
		for i := 2; i <= n; i++ {
			fmt.Fprintf(writer, "1 %d %d\n", i, avg)
		}
	}
}
