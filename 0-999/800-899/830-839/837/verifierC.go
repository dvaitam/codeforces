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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func expected(n, a, b int, x, y []int) string {
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for r1 := 0; r1 < 2; r1++ {
				w1, h1 := x[i], y[i]
				if r1 == 1 {
					w1, h1 = y[i], x[i]
				}
				for r2 := 0; r2 < 2; r2++ {
					w2, h2 := x[j], y[j]
					if r2 == 1 {
						w2, h2 = y[j], x[j]
					}
					if w1+w2 <= a && max(h1, h2) <= b {
						area := w1*h1 + w2*h2
						if area > best {
							best = area
						}
					}
					if max(w1, w2) <= a && h1+h2 <= b {
						area := w1*h1 + w2*h2
						if area > best {
							best = area
						}
					}
				}
			}
		}
	}
	return fmt.Sprintf("%d", best)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	a := rng.Intn(10) + 1
	b := rng.Intn(10) + 1
	x := make([]int, n)
	y := make([]int, n)
	for i := 0; i < n; i++ {
		x[i] = rng.Intn(10) + 1
		y[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", x[i], y[i]))
	}
	return sb.String(), expected(n, a, b, x, y)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
