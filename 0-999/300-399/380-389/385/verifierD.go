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

func expected(n int, l, r float64, px, py, a []float64) float64 {
	m := 1 << uint(n)
	dp := make([]float64, m)
	negInf := -1e300
	for i := 1; i < m; i++ {
		dp[i] = negInf
	}
	dp[0] = l

	ans := l

	for mask := 0; mask < m; mask++ {
		p := dp[mask]
		if p == negInf {
			continue
		}
		if p > ans {
			ans = p
		}
		if p >= r {
			continue
		}

		for i := 0; i < n; i++ {
			if mask&(1<<uint(i)) != 0 {
				continue
			}

			var v float64
			isA90 := math.Abs(a[i]-90.0) < 1e-9

			if isA90 {
				if p >= px[i] {
					v = r
				} else {
					v = px[i] - py[i]*py[i]/(p-px[i])
					if v > r {
						v = r
					}
				}
			} else {
				q := math.Tan(a[i] * math.Pi / 180.0)
				t := (p - px[i]) / py[i]
				den := 1.0 - t*q
				if den <= 0 {
					v = r
				} else {
					v = px[i] + py[i]*(t+q)/den
					if v > r {
						v = r
					}
				}
			}

			nmask := mask | (1 << uint(i))
			if v > dp[nmask] {
				dp[nmask] = v
			}
		}
	}

	if ans > r {
		ans = r
	}

	return ans - l
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(4) + 1
	l := rng.Intn(20) - 10
	r := l + rng.Intn(10) + 1
	px := make([]float64, n)
	py := make([]float64, n)
	ang := make([]float64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, l, r))
	for i := 0; i < n; i++ {
		xi := rng.Intn(20) - 10
		yi := rng.Intn(10) + 1
		ai := rng.Intn(89) + 1
		px[i] = float64(xi)
		py[i] = float64(yi)
		ang[i] = float64(ai)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", xi, yi, ai))
	}
	return sb.String(), expected(n, float64(l), float64(r), px, py, ang)
}

func runCase(bin, input string, expected float64) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
