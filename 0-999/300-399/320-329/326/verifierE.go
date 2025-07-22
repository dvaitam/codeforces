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

func expected(name string, n int64, h int) float64 {
	offset := int64((1 << (h + 1)) - 2)
	if name == "Alice" {
		return float64(n + offset)
	}
	return float64(n - offset)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(bin, name string, n int64, h int) error {
	input := fmt.Sprintf("%s\n%d %d\n", name, n, h)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var val float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &val); err != nil {
		return fmt.Errorf("unable to parse output: %v", err)
	}
	exp := expected(name, n, h)
	if abs(val-exp) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", exp, val)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int64, int) {
	name := "Alice"
	if rng.Intn(2) == 0 {
		name = "Bob"
	}
	return name, rng.Int63n(1000), rng.Intn(10)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		name, n, h := generateCase(rng)
		if err := runCase(bin, name, n, h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
