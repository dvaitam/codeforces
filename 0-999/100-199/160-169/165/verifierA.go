package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type point struct{ x, y int }

func expected(points []point) int {
	n := len(points)
	count := 0
	for i := 0; i < n; i++ {
		hasL, hasR, hasU, hasD := false, false, false, false
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			if points[j].x == points[i].x {
				if points[j].y > points[i].y {
					hasU = true
				}
				if points[j].y < points[i].y {
					hasD = true
				}
			}
			if points[j].y == points[i].y {
				if points[j].x > points[i].x {
					hasR = true
				}
				if points[j].x < points[i].x {
					hasL = true
				}
			}
		}
		if hasL && hasR && hasU && hasD {
			count++
		}
	}
	return count
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(20) + 1
	if rng.Float64() < 0.1 {
		n = 200
	}
	pts := make([]point, n)
	for i := 0; i < n; i++ {
		pts[i] = point{rng.Intn(21) - 10, rng.Intn(21) - 10}
	}
	var b strings.Builder
	fmt.Fprintln(&b, n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d %d\n", pts[i].x, pts[i].y)
	}
	return b.String(), expected(pts)
}

func runCase(bin, input string, exp int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
