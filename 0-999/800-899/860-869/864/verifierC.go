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

func solveC(a, b, f, k int) string {
	fuel := b
	refuels := 0
	for i := 1; i <= k; i++ {
		if i%2 == 1 { // 0->a
			if fuel < f {
				return "-1"
			}
			fuel -= f
			dist := 0
			if i == k {
				dist = a - f
			} else {
				dist = 2 * (a - f)
			}
			if b < dist {
				return "-1"
			}
			if fuel < dist {
				refuels++
				fuel = b
			}
			fuel -= a - f
		} else { // a->0
			if fuel < a-f {
				return "-1"
			}
			fuel -= a - f
			dist := 0
			if i == k {
				dist = f
			} else {
				dist = 2 * f
			}
			if b < dist {
				return "-1"
			}
			if fuel < dist {
				refuels++
				fuel = b
			}
			fuel -= f
		}
	}
	return fmt.Sprintf("%d", refuels)
}

func generateCase(rng *rand.Rand) (string, string) {
	a := rng.Intn(1000) + 1
	f := rng.Intn(a-1) + 1
	b := rng.Intn(1000) + 1
	k := rng.Intn(10) + 1
	input := fmt.Sprintf("%d %d %d %d\n", a, b, f, k)
	return input, solveC(a, b, f, k)
}

func runCandidate(bin, input string) (string, error) {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
