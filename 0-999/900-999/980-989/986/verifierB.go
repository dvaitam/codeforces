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

func solveCase(n int, perm []int) string {
	visited := make([]bool, n+1)
	cycles := 0
	for i := 1; i <= n; i++ {
		if !visited[i] {
			cycles++
			for j := i; !visited[j]; j = perm[j] {
				visited[j] = true
			}
		}
	}
	parity := (n - cycles) % 2
	petrParity := n % 2
	if parity == petrParity {
		return "Petr\n"
	}
	return "Um_nik\n"
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	p0 := rng.Perm(n)
	perm := make([]int, n+1)
	for i := 0; i < n; i++ {
		perm[i+1] = p0[i] + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", perm[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solveCase(n, perm)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
