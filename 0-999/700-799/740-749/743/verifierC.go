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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solveCaseC(n int64) string {
	if n == 1 {
		return "-1"
	}
	if n%2 == 0 {
		k := n / 2
		x := k + 2
		y := k * (k + 1)
		z := (k + 1) * (k + 2)
		return fmt.Sprintf("%d %d %d", x, y, z)
	}
	half := n/2 + 1
	maxX := half + 2000
	if maxX > 2*n {
		maxX = 2 * n
	}
	for x := half; x <= maxX; x++ {
		num := 2*x - n
		den := n * x
		g := gcd(num, den)
		a := num / g
		b := den / g
		y0 := b/a + 1
		if y0 <= x {
			y0 = x + 1
		}
		for y := y0; y < y0+10; y++ {
			d := a*y - b
			if d <= 0 {
				continue
			}
			by := b * y
			if by%d != 0 {
				continue
			}
			z := by / d
			if z <= 0 || z > 1000000000 {
				continue
			}
			if z == y || z == x {
				continue
			}
			return fmt.Sprintf("%d %d %d", x, y, z)
		}
	}
	return "-1"
}

func generateCaseC(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1000) + 1
	input := fmt.Sprintf("%d\n", n)
	expect := solveCaseC(n)
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseC(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
