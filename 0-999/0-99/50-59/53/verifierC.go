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

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1
	input := fmt.Sprintf("%d\n", n)
	return input, n
}

func validPermutation(n int, arr []int) bool {
	seen := make([]bool, n)
	diff := make(map[int]bool)
	if len(arr) != n {
		return false
	}
	for i, v := range arr {
		if v < 1 || v > n || seen[v-1] {
			return false
		}
		seen[v-1] = true
		if i > 0 {
			d := v - arr[i-1]
			if d < 0 {
				d = -d
			}
			if diff[d] {
				return false
			}
			diff[d] = true
		}
	}
	return true
}

func runCase(exe, input string, n int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	arr := make([]int, len(fields))
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		arr[i] = v
	}
	if !validPermutation(n, arr) {
		return fmt.Errorf("output is not a valid permutation")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, n := generateCase(rng)
		if err := runCase(exe, in, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
