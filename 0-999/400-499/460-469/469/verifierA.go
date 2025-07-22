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

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	p := rng.Intn(n + 1)
	q := rng.Intn(n + 1)

	perm := rng.Perm(n)
	x := perm[:p]
	perm2 := rng.Perm(n)
	y := perm2[:q]

	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	fmt.Fprintf(&sb, "%d", p)
	for _, v := range x {
		fmt.Fprintf(&sb, " %d", v+1)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d", q)
	for _, v := range y {
		fmt.Fprintf(&sb, " %d", v+1)
	}
	sb.WriteByte('\n')

	seen := make([]bool, n+1)
	for _, v := range x {
		seen[v+1] = true
	}
	for _, v := range y {
		seen[v+1] = true
	}
	ok := true
	for i := 1; i <= n; i++ {
		if !seen[i] {
			ok = false
			break
		}
	}
	exp := "I become the guy."
	if !ok {
		exp = "Oh, my keyboard!"
	}
	return sb.String(), exp
}

func runCase(bin string, input string, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
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
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
