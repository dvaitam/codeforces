package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	return a / gcd(a, b) * b
}

func verifyCase(bin string, x int64) error {
	input := fmt.Sprintf("1\n%d\n", x)
	out, err := run(bin, input)
	if err != nil {
		return err
	}
	fields := strings.Fields(out)
	if len(fields) != 2 {
		return fmt.Errorf("expected two numbers, got %q", out)
	}
	a, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q", fields[0])
	}
	b, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid number %q", fields[1])
	}
	if gcd(a, b)+lcm(a, b) != x {
		return fmt.Errorf("GCD+LCM != x for x=%d a=%d b=%d", x, a, b)
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
		x := rng.Int63n(1_000_000_000-1) + 2
		if err := verifyCase(bin, x); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
