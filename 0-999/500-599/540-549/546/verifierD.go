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

const maxLimit = 500000

var spf []int
var bigOmega []int
var pref []int64

func precompute() {
	spf = make([]int, maxLimit+1)
	for i := 2; i <= maxLimit; i++ {
		if spf[i] == 0 {
			for j := i; j <= maxLimit; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	bigOmega = make([]int, maxLimit+1)
	pref = make([]int64, maxLimit+1)
	for i := 2; i <= maxLimit; i++ {
		p := spf[i]
		bigOmega[i] = bigOmega[i/p] + 1
		pref[i] = pref[i-1] + int64(bigOmega[i])
	}
}

func solve(a, b int) string {
	return fmt.Sprintf("%d", pref[a]-pref[b])
}

func generateCase(rng *rand.Rand) (string, string) {
	a := rng.Intn(maxLimit-1) + 1
	b := rng.Intn(a) // ensures b < a
	input := fmt.Sprintf("1\n%d %d\n", a, b)
	return input, solve(a, b)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	precompute()
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
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
