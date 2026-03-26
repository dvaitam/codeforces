package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ext struct {
	m float64
	e int32
}

func solveCase(n int, X int, pr []int) string {
	c := float64(X) / 1000000.0

	lfact := make([]float64, n+1)
	for i := 1; i <= n; i++ {
		lfact[i] = lfact[i-1] + math.Log2(float64(i))
	}

	p := make([]ext, n+1)
	for r := 0; r <= n; r++ {
		if pr[r] == 0 {
			p[r] = ext{0.0, 0}
		} else {
			log2Comb := lfact[n] - lfact[r] - lfact[n-r]
			L := math.Log2(float64(pr[r])/1000000.0) - log2Comb
			e := int32(math.Floor(L)) + 1
			m := math.Exp2(L - float64(e))
			p[r] = ext{m, e}
		}
	}

	var pow2 [60]float64
	for i := 0; i < 60; i++ {
		pow2[i] = math.Exp2(float64(-i))
	}

	V := make([]float64, n+1)
	nextV := make([]float64, n+1)
	nextPArr := make([]ext, n+1)

	for d := n - 1; d >= 0; d-- {
		for r := 0; r <= d; r++ {
			a := p[r]
			b := p[r+1]

			var nextP ext
			if a.m == 0 {
				nextP = b
			} else if b.m == 0 {
				nextP = a
			} else {
				var diffE int32
				if a.e >= b.e {
					diffE = a.e - b.e
					if diffE > 59 {
						nextP = a
					} else {
						m := a.m + b.m*pow2[diffE]
						if m >= 1.0 {
							nextP = ext{m * 0.5, a.e + 1}
						} else {
							nextP = ext{m, a.e}
						}
					}
				} else {
					diffE = b.e - a.e
					if diffE > 59 {
						nextP = b
					} else {
						m := b.m + a.m*pow2[diffE]
						if m >= 1.0 {
							nextP = ext{m * 0.5, b.e + 1}
						} else {
							nextP = ext{m, b.e}
						}
					}
				}
			}

			if nextP.m == 0 {
				nextV[r] = 0.0
			} else {
				var probRed float64
				if b.m == 0 {
					probRed = 0.0
				} else {
					diff := nextP.e - b.e
					if diff > 59 {
						probRed = 0.0
					} else {
						probRed = (b.m / nextP.m) * pow2[diff]
					}
				}
				if probRed > 1.0 {
					probRed = 1.0
				}

				expected := -c + probRed*(1.0+V[r+1]) + (1.0-probRed)*V[r]
				if expected > 0.0 {
					nextV[r] = expected
				} else {
					nextV[r] = 0.0
				}
			}
			nextPArr[r] = nextP
		}

		p, nextPArr = nextPArr, p
		V, nextV = nextV, V
	}

	return fmt.Sprintf("%.12f", V[0])
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	X := rng.Intn(1000001)
	weights := make([]int, n+1)
	total := 0
	for i := 0; i <= n; i++ {
		weights[i] = rng.Intn(100) + 1
		total += weights[i]
	}
	p := make([]int, n+1)
	sum := 0
	for i := 0; i < n; i++ {
		p[i] = weights[i] * 1000000 / total
		sum += p[i]
	}
	p[n] = 1000000 - sum
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, X))
	for i := 0; i <= n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", p[i]))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, input string) error {
	// Parse input to get expected
	var n, X int
	r := strings.NewReader(input)
	fmt.Fscan(r, &n, &X)
	pr := make([]int, n+1)
	for i := 0; i <= n; i++ {
		fmt.Fscan(r, &pr[i])
	}
	expected := solveCase(n, X, pr)

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())

	var gotVal, expVal float64
	fmt.Sscan(got, &gotVal)
	fmt.Sscan(expected, &expVal)
	if math.Abs(gotVal-expVal) > 1e-4 {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
