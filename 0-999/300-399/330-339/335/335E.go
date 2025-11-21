package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func expectedBob(n, h int) float64 {
	pow2 := make([]float64, h+1)
	for i := 0; i <= h; i++ {
		pow2[i] = math.Exp2(float64(i))
	}

	pVal := make([]float64, h+1)
	for i := 0; i <= h; i++ {
		if i < h {
			pVal[i] = math.Exp2(float64(-(i + 1)))
		} else {
			pVal[i] = math.Exp2(float64(-h))
		}
	}

	cum := make([]float64, h+1)
	acc := 0.0
	for i := 0; i <= h; i++ {
		acc += pVal[i]
		cum[i] = acc
	}

	dpColumn := make([][]float64, h)
	var finalDP [][]float64

	for c := 0; c <= h; c++ {
		size := c + 1
		dpCurr := make([][]float64, n+1)
		for i := range dpCurr {
			dpCurr[i] = make([]float64, size)
		}

		norm := cum[c]
		pC := make([]float64, size)
		for t := 0; t <= c; t++ {
			pC[t] = pVal[t] / norm
		}

		beta := make([]float64, size)
		for a := 0; a <= c; a++ {
			lo := 0.0
			if a > 0 {
				lo = cum[a-1]
			}
			q := (norm - lo) / norm
			beta[a] = 1 - q
		}

		A := make([]float64, c)
		B := make([]float64, c)
		for v := 0; v < c; v++ {
			if v == 0 {
				A[v] = 0
			} else {
				A[v] = cum[v-1] / norm
			}
			B[v] = cum[v] / norm
		}

		g := make([][]float64, size)
		for a := 0; a < size; a++ {
			g[a] = make([]float64, size)
		}

		powBeta := make([]float64, size)
		for a := 0; a < size; a++ {
			powBeta[a] = 1
		}

		H := make([]float64, c)
		Q := make([]float64, c)
		powB := make([]float64, c)
		for v := 0; v < c; v++ {
			powB[v] = 1
		}

		for v := 0; v < c; v++ {
			col := dpColumn[v]
			next := 0.0
			if len(col) > 0 {
				next = col[1]
			}
			H[v] = powB[v]*next + A[v]*H[v]
			Q[v] = powB[v] + A[v]*Q[v]
			powB[v] *= B[v]
		}

		for m := 2; m <= n; m++ {
			for a := 0; a < size; a++ {
				powBeta[a] *= beta[a]
			}
			for a := 0; a < size; a++ {
				const1 := pow2[a] * (1 - powBeta[a])
				term1 := 0.0
				for t := a; t < size; t++ {
					term1 += pC[t] * g[a][t]
				}
				const2 := 0.0
				term2 := 0.0
				for v := 0; v < a; v++ {
					const2 += pow2[v] * pC[v] * Q[v]
					term2 += pC[v] * H[v]
				}
				dpCurr[m][a] = const1 + term1 + const2 + term2
			}
			for a := 0; a < size; a++ {
				for t := a; t < size; t++ {
					g[a][t] = dpCurr[m][t] + beta[a]*g[a][t]
				}
			}
			for v := 0; v < c; v++ {
				col := dpColumn[v]
				next := 0.0
				if len(col) > 0 {
					next = col[m]
				}
				H[v] = powB[v]*next + A[v]*H[v]
				Q[v] = powB[v] + A[v]*Q[v]
				powB[v] *= B[v]
			}
		}

		if c == h {
			finalDP = dpCurr
		} else {
			col := make([]float64, n+1)
			for m := 0; m <= n; m++ {
				col[m] = dpCurr[m][c]
			}
			dpColumn[c] = col
		}
	}

	ans := 1.0
	for a := 0; a <= h; a++ {
		ans += pVal[a] * finalDP[n][a]
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var who string
	var n, h int
	fmt.Fscan(in, &who)
	fmt.Fscan(in, &n, &h)

	if who == "Bob" {
		fmt.Fprintf(out, "%.10f\n", float64(n))
		return
	}

	ans := expectedBob(n, h)
	fmt.Fprintf(out, "%.10f\n", ans)
}
