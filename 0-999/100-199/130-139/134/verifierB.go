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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func steps(n int) int {
	if n == 1 {
		return 0
	}
	const inf = int(1e9)
	minSteps := inf
	for b := 1; b < n; b++ {
		x, y := n, b
		sum := 0
		for y != 0 {
			sum += x / y
			x, y = y, x%y
		}
		if x != 1 {
			continue
		}
		if sum-1 < minSteps {
			minSteps = sum - 1
		}
	}
	return minSteps
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(100000) + 1
	expect := steps(n)
	return fmt.Sprintf("%d\n", n), fmt.Sprint(expect)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
