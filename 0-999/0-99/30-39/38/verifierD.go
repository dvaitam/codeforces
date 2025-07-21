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

func solveCase(n int, x1, y1, x2, y2 []float64) int {
	minX := make([]float64, n+1)
	maxX := make([]float64, n+1)
	minY := make([]float64, n+1)
	maxY := make([]float64, n+1)
	cx := make([]float64, n+1)
	cy := make([]float64, n+1)
	w := make([]float64, n+1)
	for i := 1; i <= n; i++ {
		minX[i] = math.Min(x1[i], x2[i])
		maxX[i] = math.Max(x1[i], x2[i])
		minY[i] = math.Min(y1[i], y2[i])
		maxY[i] = math.Max(y1[i], y2[i])
		cx[i] = (x1[i] + x2[i]) / 2.0
		cy[i] = (y1[i] + y2[i]) / 2.0
		a := maxX[i] - minX[i]
		w[i] = a * a * a
	}
	const eps = 1e-9
	for k := 2; k <= n; k++ {
		for i := k; i >= 2; i-- {
			var sumW, sumX, sumY float64
			for j := i; j <= k; j++ {
				sumW += w[j]
				sumX += w[j] * cx[j]
				sumY += w[j] * cy[j]
			}
			cmX := sumX / sumW
			cmY := sumY / sumW
			sx := math.Max(minX[i-1], minX[i])
			tx := math.Min(maxX[i-1], maxX[i])
			sy := math.Max(minY[i-1], minY[i])
			ty := math.Min(maxY[i-1], maxY[i])
			if cmX < sx-eps || cmX > tx+eps || cmY < sy-eps || cmY > ty+eps {
				return k - 1
			}
		}
	}
	return n
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	x1 := make([]float64, n+1)
	y1 := make([]float64, n+1)
	x2 := make([]float64, n+1)
	y2 := make([]float64, n+1)
	// base brick
	size := float64(rng.Intn(4) + 2)
	x1[1] = 0
	y1[1] = 0
	x2[1] = size
	y2[1] = size
	for i := 2; i <= n; i++ {
		s := float64(rng.Intn(4) + 1)
		px := x1[i-1] + rng.Float64()*(x2[i-1]-x1[i-1])
		py := y1[i-1] + rng.Float64()*(y2[i-1]-y1[i-1])
		x1[i] = px - rng.Float64()*s
		y1[i] = py - rng.Float64()*s
		x2[i] = x1[i] + s
		y2[i] = y1[i] + s
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", int(x1[i]), int(y1[i]), int(x2[i]), int(y2[i])))
	}
	ans := solveCase(n, x1, y1, x2, y2)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
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
