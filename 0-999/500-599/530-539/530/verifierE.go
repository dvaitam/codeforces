package main

import (
	"bytes"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func check(nums []int64, d int64) bool {
	sum := big.NewInt(0)
	prod := big.NewInt(1)
	for _, v := range nums {
		if v <= 0 || v > 1_000_000 {
			return false
		}
		sum.Add(sum, big.NewInt(v))
		prod.Mul(prod, big.NewInt(v))
	}
	diff := big.NewInt(0).Sub(prod, sum)
	return diff.Cmp(big.NewInt(d)) == 0
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		n := rng.Intn(9) + 2
		d := rng.Intn(100)
		input := fmt.Sprintf("%d %d\n", n, d)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != n {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers, got %d\n", i+1, n, len(fields))
			os.Exit(1)
		}
		nums := make([]int64, n)
		prev := int64(-1)
		for j, f := range fields {
			v, err := strconv.ParseInt(f, 10, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: invalid integer\n", i+1)
				os.Exit(1)
			}
			if j > 0 && v < prev {
				fmt.Fprintf(os.Stderr, "case %d failed: numbers not non-decreasing\n", i+1)
				os.Exit(1)
			}
			prev = v
			nums[j] = v
		}
		if !check(nums, int64(d)) {
			fmt.Fprintf(os.Stderr, "case %d failed: numbers do not satisfy condition\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
