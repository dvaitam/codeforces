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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type heater struct {
	x, y, r int
}

func covered(x, y int, hs []heater) bool {
	for _, h := range hs {
		dx := x - h.x
		dy := y - h.y
		if dx*dx+dy*dy <= h.r*h.r {
			return true
		}
	}
	return false
}

func expected(xa, ya, xb, yb int, hs []heater) string {
	minX, maxX := xa, xb
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	minY, maxY := ya, yb
	if minY > maxY {
		minY, maxY = maxY, minY
	}
	blankets := 0
	for x := minX; x <= maxX; x++ {
		if !covered(x, minY, hs) {
			blankets++
		}
		if minY != maxY && !covered(x, maxY, hs) {
			blankets++
		}
	}
	for y := minY + 1; y <= maxY-1; y++ {
		if !covered(minX, y, hs) {
			blankets++
		}
		if minX != maxX && !covered(maxX, y, hs) {
			blankets++
		}
	}
	return fmt.Sprintf("%d", blankets)
}

func genCase(rng *rand.Rand) (string, string) {
	xa := rng.Intn(21) - 10
	ya := rng.Intn(21) - 10
	xb := rng.Intn(21) - 10
	for xb == xa {
		xb = rng.Intn(21) - 10
	}
	yb := rng.Intn(21) - 10
	for yb == ya {
		yb = rng.Intn(21) - 10
	}
	n := rng.Intn(5) + 1
	hs := make([]heater, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", xa, ya, xb, yb))
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		hs[i].x = rng.Intn(21) - 10
		hs[i].y = rng.Intn(21) - 10
		hs[i].r = rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", hs[i].x, hs[i].y, hs[i].r))
	}
	exp := expected(xa, ya, xb, yb, hs)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := run(bin, in)
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
