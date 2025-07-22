package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func expected(n int) int64 {
	var count int64
	for m := 2; m*m+1 <= n; m++ {
		for k := 1; k < m; k++ {
			if (m-k)%2 == 1 && gcd(m, k) == 1 {
				c0 := m*m + k*k
				if c0 > n {
					continue
				}
				count += int64(n / c0)
			}
		}
	}
	return count
}

func runCase(exe string, n int) error {
	input := fmt.Sprintf("%d\n", n)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	resStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(resStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", resStr)
	}
	exp := expected(n)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	tests := []int{}
	for i := 1; i <= 100; i++ { // n from 1 to 100
		tests = append(tests, i)
	}
	// some larger cases
	extras := []int{150, 200, 300, 400, 500, 750, 1000, 2000, 5000, 10000}
	tests = append(tests, extras...)
	for idx, n := range tests {
		if err := runCase(exe, n); err != nil {
			fmt.Printf("test %d (n=%d) failed: %v\n", idx+1, n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
