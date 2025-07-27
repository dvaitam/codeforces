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

type testCaseC struct {
	n int
}

func buildInputC(n int) string {
	return fmt.Sprintf("1\n%d\n", n)
}

func simulateC(n int, ops [][2]int) (int, error) {
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = i + 1
	}
	for _, p := range ops {
		a, b := p[0], p[1]
		// find indices
		idxA, idxB := -1, -1
		for i, v := range nums {
			if v == a && idxA == -1 {
				idxA = i
			} else if v == b && idxB == -1 {
				idxB = i
			}
		}
		if idxA == -1 || idxB == -1 || idxA == idxB {
			return 0, fmt.Errorf("invalid pair %d %d", a, b)
		}
		if idxA > idxB {
			idxA, idxB = idxB, idxA
		}
		// remove higher index first
		nums = append(nums[:idxB], nums[idxB+1:]...)
		nums = append(nums[:idxA], nums[idxA+1:]...)
		// insert ceil((a+b)/2)
		val := (a + b + 1) / 2
		nums = append(nums, val)
	}
	if len(nums) != 1 {
		return 0, fmt.Errorf("wrong number count %d", len(nums))
	}
	return nums[0], nil
}

func runCaseC(bin string, tc testCaseC) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputC(tc.n))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) < 1+(tc.n-1)*2 {
		return fmt.Errorf("not enough numbers")
	}
	var first int
	if _, err := fmt.Sscan(fields[0], &first); err != nil || first != 2 {
		return fmt.Errorf("first number should be 2")
	}
	ops := make([][2]int, tc.n-1)
	idx := 1
	for i := 0; i < tc.n-1; i++ {
		if idx+1 >= len(fields) {
			return fmt.Errorf("not enough pairs")
		}
		var a, b int
		if _, err := fmt.Sscan(fields[idx], &a); err != nil {
			return fmt.Errorf("bad pair")
		}
		if _, err := fmt.Sscan(fields[idx+1], &b); err != nil {
			return fmt.Errorf("bad pair")
		}
		ops[i] = [2]int{a, b}
		idx += 2
	}
	if _, err := simulateC(tc.n, ops); err != nil {
		return err
	}
	return nil
}

func generateCasesC() []testCaseC {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseC, 0, 100)
	for _, n := range []int{2, 3, 4, 5, 6, 10, 50} {
		cases = append(cases, testCaseC{n})
	}
	for len(cases) < 100 {
		cases = append(cases, testCaseC{rng.Intn(50) + 2})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesC()
	for i, tc := range cases {
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", i+1, err, tc.n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
