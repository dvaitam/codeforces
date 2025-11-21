package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	var n, m int
	var aInt, bInt int
	if _, err := fmt.Fscan(in, &n, &m, &aInt, &bInt); err != nil {
		return
	}

	yA := make([]float64, n)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		yA[i] = float64(v)
	}

	yB := make([]float64, m)
	for i := 0; i < m; i++ {
		var v int
		fmt.Fscan(in, &v)
		yB[i] = float64(v)
	}

	l := make([]float64, m)
	for i := 0; i < m; i++ {
		var v int
		fmt.Fscan(in, &v)
		l[i] = float64(v)
	}

	a := float64(aInt)
	dx := float64(bInt - aInt)

	w := make([]float64, n)
	for i := 0; i < n; i++ {
		w[i] = math.Hypot(a, yA[i])
	}

	best := 0
	ansCost := math.MaxFloat64
	ansI, ansJ := 1, 1

	for j := 0; j < m; j++ {
		y := yB[j]
		bestCost := w[best] + math.Hypot(dx, y-yA[best])
		for best+1 < n {
			nextCost := w[best+1] + math.Hypot(dx, y-yA[best+1])
			if nextCost <= bestCost {
				best++
				bestCost = nextCost
			} else {
				break
			}
		}

		total := bestCost + l[j]
		if total < ansCost {
			ansCost = total
			ansI = best + 1
			ansJ = j + 1
		}
	}

	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(out, "%d %d\n", ansI, ansJ)
	out.Flush()
}
