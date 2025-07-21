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

func solveE(n int, pts [][2]float64) string {
	x1 := pts[0][0]
	y0 := pts[0][1]
	x2 := pts[1][0]
	L, R := x1, x2
	if L > R {
		L, R = R, L
	}
	const eps = 1e-9
	for i := 0; i < n; i++ {
		xi, yi := pts[i][0], pts[i][1]
		xj, yj := pts[(i+1)%n][0], pts[(i+1)%n][1]
		dx := xj - xi
		dy := yj - yi
		A := dx
		B := -dy
		if math.Abs(B) < eps {
			if A*(y0-yi) > 0 {
				return "0"
			}
			continue
		}
		num := -A * (y0 - yi)
		bound := xi + num/B
		if B > 0 {
			if bound < R {
				R = bound
			}
		} else {
			if bound > L {
				L = bound
			}
		}
		if L > R {
			return "0"
		}
	}
	start := math.Ceil(L - eps)
	end := math.Floor(R + eps)
	cnt := int(end - start + 1)
	if cnt < 0 {
		cnt = 0
	}
	return fmt.Sprintf("%d", cnt)
}

func generateCaseE(rng *rand.Rand) (string, string) {
	x1 := float64(rng.Intn(21) - 10)
	x2 := float64(rng.Intn(21) - 10)
	if x1 == x2 {
		x2 = x1 + 1
	}
	y0 := float64(rng.Intn(21) - 10)
	y1 := y0 + float64(rng.Intn(10)+1)
	pts := [][2]float64{
		{x1, y0},
		{x2, y0},
		{x2, y1},
		{x1, y1},
	}
	n := len(pts)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%.0f %.0f\n", pts[i][0], pts[i][1])
	}
	expected := solveE(n, pts)
	return sb.String(), expected
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
	outStr := strings.TrimSpace(out.String())
	if outStr != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseE(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
