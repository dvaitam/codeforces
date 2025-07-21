package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n, m, k float64) float64 {
	if n+k < m {
		return 0.0
	}
	
	ans := 1.0
	for i := 0; i <= int(k); i++ {
		ans *= (m - float64(i)) / (n + 1 + float64(i))
	}
	return 1.0 - ans
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	for i := 0; i < 100; i++ {
		n := rng.Intn(101)
		m := rng.Intn(101)
		k := rng.Intn(11)
		
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		
		expectedOut := expected(float64(n), float64(m), float64(k))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		
		gotFloat, parseErr := strconv.ParseFloat(got, 64)
		if parseErr != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output: %v\ninput:\n%s", i+1, parseErr, input)
			os.Exit(1)
		}
		
		if math.Abs(gotFloat-expectedOut) > 1e-6 {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %.6f got %.6f (diff %.6f)\ninput:\n%s", i+1, expectedOut, gotFloat, math.Abs(gotFloat-expectedOut), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}