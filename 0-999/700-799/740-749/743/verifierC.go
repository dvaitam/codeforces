package main

import (
	"bytes"
	"fmt"
	"math/big"
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

func checkAnswer(nStr, exp, got string) error {
	exp = strings.TrimSpace(exp)
	got = strings.TrimSpace(got)
	if exp == "-1" {
		if got != "-1" {
			return fmt.Errorf("expected -1 got %s", got)
		}
		return nil
	}
	if got == "-1" {
		return fmt.Errorf("expected valid answer, got -1")
	}

	var x, y, z int64
	parsed, err := fmt.Sscanf(got, "%d %d %d", &x, &y, &z)
	if err != nil || parsed != 3 {
		return fmt.Errorf("failed to parse output: %s", got)
	}

	if x <= 0 || y <= 0 || z <= 0 || x > 1000000000 || y > 1000000000 || z > 1000000000 {
		return fmt.Errorf("values out of bounds: %d %d %d", x, y, z)
	}

	if x == y || y == z || x == z {
		return fmt.Errorf("values not distinct: %d %d %d", x, y, z)
	}

	var n int64
	fmt.Sscanf(nStr, "%d", &n)

	left := new(big.Int).Mul(big.NewInt(2), big.NewInt(x))
	left.Mul(left, big.NewInt(y))
	left.Mul(left, big.NewInt(z))

	sum := new(big.Int).Mul(big.NewInt(x), big.NewInt(y))
	sum.Add(sum, new(big.Int).Mul(big.NewInt(y), big.NewInt(z)))
	sum.Add(sum, new(big.Int).Mul(big.NewInt(z), big.NewInt(x)))

	right := new(big.Int).Mul(big.NewInt(n), sum)

	if left.Cmp(right) != 0 {
		return fmt.Errorf("equation not satisfied: 2/%d != 1/%d + 1/%d + 1/%d (got %d %d %d)", n, x, y, z, x, y, z)
	}
	return nil
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
		if err := checkAnswer(in, exp, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}