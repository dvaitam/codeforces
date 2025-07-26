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

func runProg(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int64) string {
	switch {
	case n < 1000:
		return fmt.Sprintf("%d", n)
	case n < 1_000_000:
		k := (n + 500) / 1000
		if k == 1000 {
			return "1M"
		}
		return fmt.Sprintf("%dK", k)
	default:
		m := (n + 500_000) / 1_000_000
		return fmt.Sprintf("%dM", m)
	}
}

func runCase(bin string, n int64) error {
	input := fmt.Sprintf("1\n%d\n", n)
	want := expected(n)
	out, err := runProg(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	if out != want {
		return fmt.Errorf("expected %s got %s (n=%d)", want, out, n)
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
		n := rng.Int63n(2_000_000_001)
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
