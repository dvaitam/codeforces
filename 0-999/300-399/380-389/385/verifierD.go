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
	cosA := make([]float64, n)
	sinA := make([]float64, n)
	for i := 0; i < n; i++ {
		rad := a[i] * math.Pi / 180.0
		cosA[i] = math.Cos(rad)
		sinA[i] = math.Sin(rad)
	}
	size := 1 << uint(n)
	dp := make([]float64, size)
	for i := range dp {
		dp[i] = l
	}
	for mask := 0; mask < size; mask++ {
		base := dp[mask]
		for i := 0; i < n; i++ {
			if mask>>i&1 == 0 {
				dx0 := base - px[i]
				dy0 := -py[i]
				dirX := dx0*cosA[i] - dy0*sinA[i]
				dirY := dx0*sinA[i] + dy0*cosA[i]
				np := r
				if dirY < 0 {
					t := -py[i] / dirY
					np = px[i] + dirX*t
					if np > r {
						np = r
					}
				}
				next := mask | (1 << i)
				if np > dp[next] {
					dp[next] = np
				}
			}
		}
	}
	return dp[size-1] - l
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(4) + 1
	l := rng.Float64()*20 - 10
	r := l + rng.Float64()*10 + 1e-3
	px := make([]float64, n)
	py := make([]float64, n)
	ang := make([]float64, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %.6f %.6f\n", n, l, r))
	for i := 0; i < n; i++ {
		px[i] = rng.Float64()*20 - 10
		py[i] = rng.Float64()*10 + 1
		ang[i] = rng.Float64()*89 + 1
		sb.WriteString(fmt.Sprintf("%.6f %.6f %.6f\n", px[i], py[i], ang[i]))
	}
	return sb.String(), expected(n, l, r, px, py, ang)
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
