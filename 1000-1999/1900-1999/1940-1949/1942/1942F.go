package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	if _, err := fmt.Fscan(in, &n, &q); err != nil {
		return
	}

	a := make([]float64, n)
	for i := 0; i < n; i++ {
		var v float64
		fmt.Fscan(in, &v)
		a[i] = v
	}

	f := make([]float64, n+1)
	for i := 0; i < n; i++ {
		f[i+1] = math.Sqrt(f[i] + a[i])
	}

	const eps = 1e-9
	for ; q > 0; q-- {
		var k int
		var x float64
		fmt.Fscan(in, &k, &x)
		k--
		a[k] = x

		prev := f[k]
		for i := k; i < n; i++ {
			val := math.Sqrt(prev + a[i])
			diff := math.Abs(val - f[i+1])
			f[i+1] = val
			prev = val
			if diff < eps && val >= 1 {
				break
			}
		}

		fmt.Fprintln(out, int64(math.Floor(f[n])))
	}
}
