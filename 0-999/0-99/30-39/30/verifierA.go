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

func solve(a, b, n int) string {
	if a == 0 {
		if b == 0 {
			return "0"
		}
		return "No solution"
	}
	for x := -1000; x <= 1000; x++ {
		val := float64(a) * math.Pow(float64(x), float64(n))
		if math.Abs(val-float64(b)) < 0.5 {
			return fmt.Sprint(x)
		}
	}
	return "No solution"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	if rng.Float64() < 0.5 {
		// generate solvable case
		a := rng.Intn(21) - 10
		if a == 0 {
			a = 1
		}
		x := rng.Intn(11) - 5
		b := int(math.Round(float64(a) * math.Pow(float64(x), float64(n))))
		input := fmt.Sprintf("%d %d %d\n", a, b, n)
		return input, solve(a, b, n)
	}
	for {
		a := rng.Intn(41) - 20
		b := rng.Intn(41) - 20
		if a == 0 {
			if b != 0 {
				input := fmt.Sprintf("%d %d %d\n", a, b, n)
				return input, solve(a, b, n)
			}
			continue
		}
		if solve(a, b, n) == "No solution" {
			input := fmt.Sprintf("%d %d %d\n", a, b, n)
			return input, "No solution"
		}
	}
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
	if outStr != expected {
		return fmt.Errorf("expected %s got %s", expected, outStr)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
