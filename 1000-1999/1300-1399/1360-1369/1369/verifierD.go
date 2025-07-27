package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 1000000007
const MAXN = 2000000

var precomputed []int64

func initPrecompute() {
	precomputed = make([]int64, MAXN+1)
	precomputed[1] = 0
	precomputed[2] = 0
	precomputed[3] = 4
	precomputed[4] = 4
	for i := 5; i <= MAXN; i++ {
		precomputed[i] = (2 * precomputed[i-1]) % MOD
		switch i % 6 {
		case 3, 5:
			precomputed[i] = (precomputed[i] + 4) % MOD
		case 4:
			precomputed[i] = (precomputed[i] + MOD - 4) % MOD
		}
	}
}

func solveD(n int) int64 {
	return precomputed[n]
}

func runBinary(binPath string, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initPrecompute()
	rand.Seed(4)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(MAXN-1) + 1
		input := fmt.Sprintf("1\n%d\n", n)
		expected := fmt.Sprintf("%d\n", solveD(n))
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", t+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", t+1, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
