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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func expected(n int, x0, y0 int, pts [][2]int) string {
	dirs := make(map[[2]int]struct{})
	for _, p := range pts {
		dx := p[0] - x0
		dy := p[1] - y0
		g := gcd(abs(dx), abs(dy))
		dx /= g
		dy /= g
		if dx < 0 || (dx == 0 && dy < 0) {
			dx = -dx
			dy = -dy
		}
		dirs[[2]int{dx, dy}] = struct{}{}
	}
	return fmt.Sprintf("%d", len(dirs))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	x0 := rng.Intn(21) - 10
	y0 := rng.Intn(21) - 10
	pts := make([][2]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x0, y0))
	for i := 0; i < n; i++ {
		x := rng.Intn(21) - 10
		y := rng.Intn(21) - 10
		pts[i] = [2]int{x, y}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	expect := expected(n, x0, y0, pts)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// deterministic cases
	fixed := []struct {
		n      int
		x0, y0 int
		pts    [][2]int
	}{
		{1, 0, 0, [][2]int{{1, 1}}},
		{3, 0, 0, [][2]int{{1, 0}, {2, 0}, {3, 0}}},
		{3, 1, 1, [][2]int{{2, 2}, {0, 0}, {3, 1}}},
	}
	for idx, f := range fixed {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", f.n, f.x0, f.y0))
		for _, p := range f.pts {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
		}
		exp := expected(f.n, f.x0, f.y0, f.pts)
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "edge case %d failed: expected %s got %s\ninput:\n%s", idx+1, exp, out, sb.String())
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
