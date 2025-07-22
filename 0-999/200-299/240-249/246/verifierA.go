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

// runOnePassSort applies Valera's buggy algorithm
func runOnePassSort(a []int) []int {
	b := make([]int, len(a))
	copy(b, a)
	for i := 0; i+1 < len(b); i++ {
		if b[i] > b[i+1] {
			b[i], b[i+1] = b[i+1], b[i]
		}
	}
	return b
}

func isSorted(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i-1] > a[i] {
			return false
		}
	}
	return true
}

func checkOutput(n int, out string) error {
	out = strings.TrimSpace(out)
	if n < 3 {
		if out != "-1" {
			return fmt.Errorf("expected -1, got %q", out)
		}
		return nil
	}
	parts := strings.Fields(out)
	if len(parts) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(parts))
	}
	arr := make([]int, n)
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			return fmt.Errorf("not an integer: %q", p)
		}
		if v < 1 || v > 100 {
			return fmt.Errorf("value %d out of range", v)
		}
		arr[i] = v
	}
	b := runOnePassSort(arr)
	if isSorted(b) {
		return fmt.Errorf("array still sorted after algorithm")
	}
	return nil
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return checkOutput(n, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		if err := runCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", i+1, err, n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
