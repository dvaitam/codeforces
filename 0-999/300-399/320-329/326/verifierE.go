package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Embedded solver for 326E
func solve326E(input string) (string, error) {
	var name string
	var n int
	var h int
	_, err := fmt.Sscan(input, &name, &n, &h)
	if err != nil {
		return "", fmt.Errorf("parse error: %v", err)
	}

	if name == "Alice" {
		if n == 0 {
			return fmt.Sprintf("%.10f", 1.0), nil
		}
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
		powF[0] = make([]float64, n+1)
		powF[0][0] = 1.0
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
		return fmt.Sprintf("%.10f", expectedBob), nil
	} else {
		target := n - 1
		if target <= 0 {
			return fmt.Sprintf("%.10f", 0.0), nil
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

		return fmt.Sprintf("%.10f", numerator/denominator), nil
	}
}

func absf(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(bin, name string, n int64, h int) error {
	input := fmt.Sprintf("%s\n%d %d\n", name, n, h)

	expStr, err := solve326E(input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	var exp float64
	if _, err := fmt.Sscan(expStr, &exp); err != nil {
		return fmt.Errorf("oracle output parse error: %v", err)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out strings.Builder
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var val float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("unable to parse output: %v", err)
	}

	tol := math.Max(1e-6, absf(exp)*1e-9)
	if absf(val-exp) > tol {
		return fmt.Errorf("expected %.10f got %.10f", exp, val)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int64, int) {
	name := "Alice"
	if rng.Intn(2) == 0 {
		name = "Bob"
	}
	return name, rng.Int63n(1000), rng.Intn(10)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		name, n, h := generateCase(rng)
		if err := runCase(bin, name, n, h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
