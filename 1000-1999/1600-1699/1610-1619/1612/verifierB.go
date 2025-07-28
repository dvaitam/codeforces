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

func possible(n, a, b int) bool {
	half := n / 2
	return a <= half+1 && b >= half && a != b
}

func check(n, a, b int, output string) error {
	output = strings.TrimSpace(output)
	if output == "-1" {
		if possible(n, a, b) {
			return fmt.Errorf("expected valid permutation but got -1")
		}
		return nil
	}
	if !possible(n, a, b) {
		return fmt.Errorf("expected -1 but got %q", output)
	}
	parts := strings.Fields(output)
	if len(parts) != n {
		return fmt.Errorf("expected %d integers, got %d", n, len(parts))
	}
	seen := make([]bool, n+1)
	vals := make([]int, n)
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("invalid integer %q", p)
		}
		if v < 1 || v > n || seen[v] {
			return fmt.Errorf("invalid permutation value %d", v)
		}
		seen[v] = true
		vals[i] = v
	}
	half := n / 2
	minLeft := vals[0]
	for i := 1; i < half; i++ {
		if vals[i] < minLeft {
			minLeft = vals[i]
		}
	}
	maxRight := vals[half]
	for i := half + 1; i < n; i++ {
		if vals[i] > maxRight {
			maxRight = vals[i]
		}
	}
	if minLeft != a || maxRight != b {
		return fmt.Errorf("constraints not satisfied")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(50)*2 + 2 // even from 2..100
		a := rand.Intn(n) + 1
		b := rand.Intn(n-1) + 1
		if b >= a {
			b++
			if b > n {
				b = 1
			}
		}
		input := fmt.Sprintf("1\n%d %d %d\n", n, a, b)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := check(n, a, b, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: n=%d a=%d b=%d\n", i+1, err, n, a, b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
