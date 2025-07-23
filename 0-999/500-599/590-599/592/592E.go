package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	var c, d int64
	if _, err := fmt.Fscan(in, &n, &c, &d); err != nil {
		return
	}
	angles := make([]float64, n)
	const twoPi = 2 * math.Pi
	for i := 0; i < n; i++ {
		var r, w int64
		fmt.Fscan(in, &r, &w)
		x := float64(r - c)
		y := float64(w - d)
		ang := math.Atan2(y, x)
		if ang < 0 {
			ang += twoPi
		}
		angles[i] = ang
	}
	sort.Float64s(angles)
	ext := make([]float64, 2*n)
	copy(ext, angles)
	for i := 0; i < n; i++ {
		ext[i+n] = angles[i] + twoPi
	}
	var bad int64
	j := 0
	eps := 1e-12
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		for j < i+n && ext[j]-ext[i] <= math.Pi+eps {
			j++
		}
		m := j - i - 1
		if m >= 2 {
			bad += int64(m*(m-1)) / 2
		}
	}
	total := int64(n) * int64(n-1) * int64(n-2) / 6
	good := total - bad
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintln(out, good)
	out.Flush()
}
