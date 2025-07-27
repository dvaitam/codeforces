package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
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

const mod = 1000000007

func fib(n int) int {
	if n == 0 {
		return 0
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		a, b = b, (a+b)%mod
	}
	return b % mod
}

func genCases() []int {
	rng := rand.New(rand.NewSource(1306))
	cases := make([]int, 100)
	for i := range cases {
		cases[i] = rng.Intn(1000)
	}
	return cases
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	want := fib(n)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	cases := genCases()
	for i, n := range cases {
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
