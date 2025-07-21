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

func expected(n, m, a, b, c int) bool {
	if n%2 == 1 && m%2 == 1 {
		return false
	}
	
	needA := 0
	needB := 0
	
	if n%2 == 1 {
		needA = m / 2
	}
	if m%2 == 1 {
		needB = n / 2
	}
	
	if a < needA || b < needB {
		return false
	}
	
	remainingA := a - needA
	remainingB := b - needB
	evenCells := (n / 2) * (m / 2)
	
	needExtra := evenCells - c
	if needExtra < 0 {
		return true
	}
	
	return 2*needExtra <= remainingA+remainingB
}

func validateOutput(output []string, n, m int) bool {
	if output[0] == "IMPOSSIBLE" {
		return true
	}
	
	if len(output) != n {
		return false
	}
	
	for i := 0; i < n; i++ {
		if len(output[i]) != m {
			return false
		}
	}
	
	return true
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		a := rng.Intn(100)
		b := rng.Intn(100)
		c := rng.Intn(100)
		
		input := fmt.Sprintf("%d %d %d %d %d\n", n, m, a, b, c)
		
		expectedPossible := expected(n, m, a, b, c)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		
		lines := strings.Split(got, "\n")
		isPossible := lines[0] != "IMPOSSIBLE"
		
		if isPossible != expectedPossible {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %t got %t\ninput:\n%s", i+1, expectedPossible, isPossible, input)
			os.Exit(1)
		}
		
		if !validateOutput(lines, n, m) {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output format\ninput:\n%s", i+1, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}