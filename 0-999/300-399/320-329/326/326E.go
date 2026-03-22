package main

import (
	"fmt"
	"math"
)

func main() {
	var name string
	var n int
	var h int
	if _, err := fmt.Scan(&name, &n, &h); err != nil {
		return
	}

	if name == "Alice" {
		f := make([]float64, h+2)
		F := make([]float64, h+2)
		for x := 1; x <= h+1; x++ {
			if x <= h {
				f[x] = math.Pow(2, float64(-x))
				F[x] = 1.0 - math.Pow(2, float64(-x))
			} else {
				f[x] = math.Pow(2, float64(-h-1))
				F[x] = 1.0
			}
		}

		powF := make([][]float64, h+2)
		for x := 1; x <= h+1; x++ {
			powF[x] = make([]float64, n+1)
			powF[x][0] = 1.0
			for k := 1; k <= n; k++ {
				powF[x][k] = powF[x][k-1] * F[x]
			}
		}

		T1 := make([]float64, n+1)
		for m := 1; m <= n; m++ {
			sum := 0.0
			for x := 1; x <= h+1; x++ {
				inner := 0.0
				for s := 1; s <= x; s++ {
					ps := powF[s][m] - powF[s-1][m]
					inner += ps * math.Pow(2, float64(s-1))
				}
				sum += f[x] * inner
			}
			T1[m] = sum
		}

		expectedBob := 1.0
		for j := 1; j <= n-1; j++ {
			m := n - j
			k := j - 1
			term1 := T1[m]
			term2 := 0.0
			for x := 1; x <= h+1; x++ {
				term2 += f[x] * powF[x][k] * math.Pow(2, float64(x-1)) * (1.0 - powF[x][m])
			}
			expectedBob += term1 + term2
		}
		fmt.Printf("%.10f\n", expectedBob)
	} else {
		target := n - 1
		if target == 0 {
			fmt.Printf("%.10f\n", 0.0)
			return
		}

		fArray := make([]float64, target+1)
		fArray[0] = 1.0

		for x := 1; x <= h+1; x++ {
			K := 1 << (x - 1)
			if K > target {
				continue
			}
			for i := K; i <= target; i++ {
				fArray[i] += fArray[i-K] * 0.5
			}
		}

		pX := 0.0
		for x := 1; x <= h+1; x++ {
			if target == (1 << (x - 1)) {
				if x <= h {
					pX = math.Pow(2, float64(-x))
				} else {
					pX = math.Pow(2, float64(-h-1))
				}
				break
			}
		}

		denominator := fArray[target] - pX

		S := make([]float64, target+1)
		for x := 1; x <= h+1; x++ {
			K := 1 << (x - 1)
			if K > target {
				continue
			}
			var DX float64
			if x <= h {
				DX = math.Pow(2, float64(x-2))
			} else {
				DX = math.Pow(2, float64(h-1))
			}
			limit := target / K
			halfPow := 1.0
			for j := 1; j <= limit; j++ {
				S[j*K] += DX * halfPow
				halfPow *= 0.5
			}
		}

		numerator := 0.0
		for i := 0; i <= target; i++ {
			numerator += fArray[i] * S[target-i]
		}
		numerator -= pX

		fmt.Printf("%.10f\n", numerator/denominator)
	}
}