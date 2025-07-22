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

type Point struct{ x, y int64 }

func dist(a, b Point) int64 {
	dx := a.x - b.x
	dy := a.y - b.y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	return dx*dx + dy*dy
}

func solve(n, m int64) string {
	if n == 0 {
		return fmt.Sprintf("0 1\n0 %d\n0 0\n0 %d\n", m, m-1)
	}
	if m == 0 {
		return fmt.Sprintf("1 0\n%d 0\n0 0\n%d 0\n", n, n-1)
	}
	cand := []Point{{0, 0}, {n, 0}, {0, m}, {n, m}, {1, 0}, {0, 1}, {n - 1, 0}, {0, m - 1}, {n, 1}, {1, m}, {n - 1, m}, {n, m - 1}}
	uniq := make(map[Point]bool)
	points := make([]Point, 0, len(cand))
	for _, p := range cand {
		if p.x < 0 || p.x > n || p.y < 0 || p.y > m {
			continue
		}
		if !uniq[p] {
			uniq[p] = true
			points = append(points, p)
		}
	}
	var bestSeq []Point
	bestSum := int64(-1)
	var perm func([]Point, int)
	perm = func(a []Point, l int) {
		if l == len(a)-1 {
			s := dist(a[0], a[1]) + dist(a[1], a[2]) + dist(a[2], a[3])
			if s > bestSum {
				bestSum = s
				bestSeq = append([]Point(nil), a...)
			}
			return
		}
		for i := l; i < len(a); i++ {
			a[l], a[i] = a[i], a[l]
			perm(a, l+1)
			a[l], a[i] = a[i], a[l]
		}
	}
	L := len(points)
	for i := 0; i < L; i++ {
		for j := i + 1; j < L; j++ {
			for k := j + 1; k < L; k++ {
				for l := k + 1; l < L; l++ {
					subset := []Point{points[i], points[j], points[k], points[l]}
					perm(subset, 0)
				}
			}
		}
	}
	var sb strings.Builder
	for _, p := range bestSeq {
		fmt.Fprintf(&sb, "%d %d\n", p.x, p.y)
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	var n, m int64
	for {
		n = int64(rng.Intn(1001))
		m = int64(rng.Intn(1001))
		if (n+1)*(m+1) >= 4 {
			break
		}
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	expected := solve(n, m)
	return input, expected
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
		return fmt.Errorf("expected\n%s\ngot\n%s", expected, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
