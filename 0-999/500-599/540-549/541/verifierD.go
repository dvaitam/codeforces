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
	// Check if the binary is a python script or go file, otherwise execute directly
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else if strings.HasSuffix(bin, ".py") {
		cmd = exec.Command("python3", bin)
	} else {
		// Assuming compiled binary
		if !strings.Contains(bin, "/") {
			bin = "./" + bin
		}
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func J(x int64) int64 {
	sum := int64(0)
	for k := int64(1); k*k <= x; k++ {
		if x%k == 0 {
			if gcd(k, x/k) == 1 {
				sum += k
			}
			d := x / k
			if d != k && gcd(d, x/d) == 1 {
				sum += d
			}
		}
	}
	return sum
}

func expected(A int64) string {
	cnt := int64(0)
	// Since J(x) >= x for all x >= 1 (equality only at x=1),
	// if J(x) = A, then x <= A.
	for x := int64(1); x <= A; x++ {
		if J(x) == A {
			cnt++
		}
	}
	return fmt.Sprintf("%d", cnt)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	
		// Add the regression test case
		regA := int64(2520)
	regInput := fmt.Sprintf("%d\n", regA)
	regExp := expected(regA) // Should be correct now
	regGot, err := run(bin, regInput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "regression test failed: %v\n", err)
		os.Exit(1)
	}
	if strings.TrimSpace(regGot) != strings.TrimSpace(regExp) {
		fmt.Fprintf(os.Stderr, "regression test failed: expected %s got %s\ninput:\n%s", regExp, regGot, regInput)
		os.Exit(1)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		var A int64
		// Generate small A to keep brute force verification feasible
		if rng.Intn(3) == 0 {
			A = int64(rng.Intn(2000) + 1)
		} else {
			x := int64(rng.Intn(2000) + 1)
			A = J(x)
		}
		
		input := fmt.Sprintf("%d\n", A)
		exp := expected(A)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}