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

type point struct{ x, y int }

func solveCase(n, k int, pts []point) float64 {
	var d float64
	for i := 1; i < n; i++ {
		dx := float64(pts[i-1].x - pts[i].x)
		dy := float64(pts[i-1].y - pts[i].y)
		d += math.Hypot(dx, dy)
	}
	return d * float64(k) / 50.0
}

func generateCase(rng *rand.Rand) (string, float64) {
	n := rng.Intn(99) + 2
	k := rng.Intn(1000) + 1
	pts := make([]point, n)
	used := make(map[[2]int]bool)
	for i := range pts {
		for {
			x := rng.Intn(41) - 20
			y := rng.Intn(41) - 20
			if !used[[2]int{x, y}] {
				used[[2]int{x, y}] = true
				pts[i] = point{x, y}
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, k, pts[0].x, pts[0].y))
	for i := 1; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", pts[i].x, pts[i].y))
	}
	return sb.String(), solveCase(n, k, pts)
}

func runCase(bin, input string, expected float64) error {
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
	var got float64
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if math.Abs(got-expected) > 1e-4*math.Max(1.0, math.Abs(expected)) {
		return fmt.Errorf("expected %.6f got %.6f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
