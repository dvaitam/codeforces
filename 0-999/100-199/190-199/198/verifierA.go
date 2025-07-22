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

func solve(k, b, n, t int64) int64 {
	cur := int64(1)
	steps := int64(0)
	for steps < n && k*cur+b <= t {
		cur = k*cur + b
		steps++
	}
	return n - steps
}

func generateCase(rng *rand.Rand) (string, string) {
	k := int64(rng.Intn(5) + 1)
	b := int64(rng.Intn(5))
	n := int64(rng.Intn(10) + 1)
	tVal := int64(rng.Intn(1000))
	input := fmt.Sprintf("%d %d %d %d\n", k, b, n, tVal)
	ans := solve(k, b, n, tVal)
	return input, fmt.Sprintf("%d", ans)
}

func runCase(bin string, input, expected string) error {
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
		return fmt.Errorf("expected %s got %s", expected, got)
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
