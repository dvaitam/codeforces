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

type Point struct{ x, y int }

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func expected(points []Point) int64 {
	n := len(points)
	if n < 3 {
		return 0
	}
	total := int64(n) * int64(n-1) * int64(n-2) / 6
	var deg int64
	for i := 0; i < n; i++ {
		mp := make(map[[2]int]int)
		xi, yi := points[i].x, points[i].y
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			dx := points[j].x - xi
			dy := points[j].y - yi
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			mp[[2]int{dx, dy}]++
		}
		for _, c := range mp {
			if c >= 2 {
				deg += int64(c * (c - 1) / 2)
			}
		}
	}
	deg /= 3
	return total - deg
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	pts := make([]Point, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		x := rng.Intn(201) - 100
		y := rng.Intn(201) - 100
		pts[i] = Point{x, y}
		fmt.Fprintf(&sb, "%d %d\n", x, y)
	}
	exp := expected(pts)
	return sb.String(), fmt.Sprintf("%d", exp)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
