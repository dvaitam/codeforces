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

type line struct {
	a, b, c int64
}

type point struct {
	x, y, s int64
}

func abs(v int64) int64 {
	if v < 0 {
		return -v
	}
	return v
}

func gcd(a, b int64) int64 {
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
		return -a
	}
	return a
}

func solveCase(n int, lines []line) int64 {
	pts := make([]point, n)
	for i, l := range lines {
		pts[i] = point{x: l.a * l.c, y: l.b * l.c, s: l.a*l.a + l.b*l.b}
	}
	var ans int64
	for i := 0; i < n; i++ {
		xi, yi, si := pts[i].x, pts[i].y, pts[i].s
		dup := 0
		slopes := make(map[[2]int64]int)
		for j := i + 1; j < n; j++ {
			xj, yj, sj := pts[j].x, pts[j].y, pts[j].s
			dx := xj*si - xi*sj
			dy := yj*si - yi*sj
			if dx == 0 && dy == 0 {
				dup++
				continue
			}
			g := gcd(abs(dx), abs(dy))
			dx /= g
			dy /= g
			if dx < 0 || (dx == 0 && dy < 0) {
				dx = -dx
				dy = -dy
			}
			slopes[[2]int64{dx, dy}]++
		}
		for _, m := range slopes {
			ans += int64(m*(m-1)/2) + int64(dup*m)
		}
		ans += int64(dup * (dup - 1) / 2)
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 3
	lines := make([]line, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		a := int64(rng.Intn(11) - 5)
		b := int64(rng.Intn(11) - 5)
		for a == 0 && b == 0 {
			a = int64(rng.Intn(11) - 5)
			b = int64(rng.Intn(11) - 5)
		}
		c := int64(rng.Intn(11) - 5)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
		lines[i] = line{a: a, b: b, c: c}
	}
	expected := fmt.Sprintf("%d", solveCase(n, lines))
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
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, stderr.String())
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
