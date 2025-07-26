package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func extendedGCD(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := extendedGCD(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

func ceilDiv(a, b int64) int64 {
	if b < 0 {
		a, b = -a, -b
	}
	if a >= 0 {
		return (a + b - 1) / b
	}
	return a / b
}

func floorDiv(a, b int64) int64 {
	if b < 0 {
		a, b = -a, -b
	}
	if a >= 0 {
		return a / b
	}
	return -((-a + b - 1) / b)
}

func expected(a, h, w int64) (float64, bool) {
	if h < a || w < a {
		return 0, false
	}
	A := h + a
	B := w + a
	C := w - h
	g, x0, y0 := extendedGCD(A, B)
	if C%g != 0 {
		return 0, false
	}
	factor := C / g
	m0 := x0 * factor
	n0 := -y0 * factor
	Bdiv := B / g
	Adiv := A / g
	lowN := ceilDiv(1-n0, Adiv)
	highN := floorDiv(h/a-n0, Adiv)
	lowM := ceilDiv(1-m0, Bdiv)
	highM := floorDiv(w/a-m0, Bdiv)
	low := lowN
	if lowM > low {
		low = lowM
	}
	high := highN
	if highM < high {
		high = highM
	}
	if low > high {
		return 0, false
	}
	t := high
	n := n0 + Adiv*t
	x := float64(h-n*a) / float64(n+1)
	return x, true
}

func run(bin, input string) (string, error) {
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

func genCase(rng *rand.Rand) (string, float64, bool) {
	a := int64(rng.Intn(10) + 1)
	h := int64(rng.Intn(20) + int(a))
	w := int64(rng.Intn(20) + int(a))
	if rng.Float64() < 0.1 {
		h = int64(rng.Intn(50) + int(a))
		w = int64(rng.Intn(50) + int(a))
	}
	input := fmt.Sprintf("%d %d %d\n", a, h, w)
	val, ok := expected(a, h, w)
	if !ok {
		return input, 0, false
	}
	return input, val, true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp, ok := genCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if !ok {
			if strings.TrimSpace(out) != "-1" {
				fmt.Fprintf(os.Stderr, "case %d failed: expected -1 got %s\ninput:\n%s", i+1, strings.TrimSpace(out), input)
				os.Exit(1)
			}
			continue
		}
		got, err := strconv.ParseFloat(strings.TrimSpace(out), 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if diff := got - exp; diff < -1e-6 || diff > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
