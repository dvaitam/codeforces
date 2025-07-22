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

// compute digital root of a decimal string
func digitalRoot(s string) int {
	if s == "0" {
		return 0
	}
	sum := 0
	for _, ch := range s {
		sum += int(ch - '0')
	}
	for sum >= 10 {
		tmp := 0
		for sum > 0 {
			tmp += sum % 10
			sum /= 10
		}
		sum = tmp
	}
	return sum
}

func check(k, d int, out string) error {
	out = strings.TrimSpace(out)
	if out == "No solution" {
		// Only no solution case is when k>=2 and d==0
		if k >= 2 && d == 0 {
			return nil
		}
		return fmt.Errorf("solution exists but got 'No solution'")
	}
	if len(out) != k {
		return fmt.Errorf("expected length %d got %d", k, len(out))
	}
	if out[0] == '0' && !(k == 1 && out == "0") {
		return fmt.Errorf("leading zero")
	}
	for _, ch := range out {
		if ch < '0' || ch > '9' {
			return fmt.Errorf("non-digit character")
		}
	}
	dr := digitalRoot(out)
	if dr != d {
		return fmt.Errorf("digital root %d != %d", dr, d)
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int, int) {
	k := rng.Intn(1000) + 1
	d := rng.Intn(10)
	return fmt.Sprintf("%d %d\n", k, d), k, d
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, k, d := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := check(k, d, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
