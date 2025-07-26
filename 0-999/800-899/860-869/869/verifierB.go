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

func solveB(a, b uint64) uint64 {
	if b-a >= 5 {
		return 0
	}
	res := uint64(1)
	for i := a + 1; i <= b; i++ {
		res = (res * (i % 10)) % 10
	}
	return res % 10
}

func run(bin string, input string) (string, error) {
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const max = uint64(1e18)
	for i := 0; i < 100; i++ {
		a := rng.Uint64() % (max - 10)
		diff := rng.Uint64() % 10
		b := a + diff
		if b > max {
			b = max
		}
		input := fmt.Sprintf("%d %d\n", a, b)
		expected := fmt.Sprintf("%d", solveB(a, b))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
