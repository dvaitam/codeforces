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

func solveCaseD(n int, c [][]int64) int64 {
	size := 2 * n
	var sum int64
	for i := n; i < size; i++ {
		for j := n; j < size; j++ {
			sum += c[i][j]
		}
	}
	candidates := []int64{
		c[0][n], c[0][size-1], c[n-1][n], c[n-1][size-1],
		c[n][0], c[n][n-1], c[size-1][0], c[size-1][n-1],
	}
	best := candidates[0]
	for _, v := range candidates[1:] {
		if v < best {
			best = v
		}
	}
	return sum + best
}

func runCandidate(bin, input string) (string, error) {
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

func generateCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 1
	size := 2 * n
	c := make([][]int64, size)
	for i := 0; i < size; i++ {
		c[i] = make([]int64, size)
		for j := 0; j < size; j++ {
			c[i][j] = int64(rng.Intn(10))
		}
	}
	input := fmt.Sprintf("1\n%d\n", n)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", c[i][j])
		}
		input += "\n"
	}
	exp := fmt.Sprintf("%d", solveCaseD(n, c))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseD(rng)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
