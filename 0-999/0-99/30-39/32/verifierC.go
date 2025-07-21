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

func expected(n, m, s int64) int64 {
	qx := n / s
	rx := n % s
	var sumX int64
	if rx > 0 {
		sumX = rx * (qx + 1)
	} else {
		sumX = s * qx
	}
	qy := m / s
	ry := m % s
	var sumY int64
	if ry > 0 {
		sumY = ry * (qy + 1)
	} else {
		sumY = s * qy
	}
	return sumX * sumY
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Int63n(1000) + 1
	m := rng.Int63n(1000) + 1
	s := rng.Int63n(1000) + 1
	input := fmt.Sprintf("%d %d %d\n", n, m, s)
	exp := fmt.Sprintf("%d", expected(n, m, s))
	return input, exp
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
		out, err := run(bin, in)
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
