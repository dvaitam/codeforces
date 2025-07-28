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

type Point struct {
	x int64
	y int64
}

func rectArea(pts []Point) int64 {
	if len(pts) == 0 {
		return 0
	}
	minX, maxX := pts[0].x, pts[0].x
	minY, maxY := pts[0].y, pts[0].y
	for _, p := range pts[1:] {
		if p.x < minX {
			minX = p.x
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.y > maxY {
			maxY = p.y
		}
	}
	return (maxX - minX) * (maxY - minY)
}

func solveCase(points []Point) int64 {
	n := len(points)
	best := int64(math.MaxInt64)
	for mask := 0; mask < (1 << n); mask++ {
		var a, b []Point
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				a = append(a, points[i])
			} else {
				b = append(b, points[i])
			}
		}
		area := rectArea(a) + rectArea(b)
		if area < best {
			best = area
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	pts := make([]Point, n)
	for i := 0; i < n; i++ {
		pts[i] = Point{int64(rng.Intn(11)), int64(rng.Intn(11))}
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i := 0; i < n; i++ {
		input += fmt.Sprintf("%d %d\n", pts[i].x, pts[i].y)
	}
	ans := solveCase(pts)
	return input, fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
