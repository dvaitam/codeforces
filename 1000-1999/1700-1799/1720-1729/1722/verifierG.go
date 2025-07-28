package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseG struct {
	n int
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseG {
	rand.Seed(48)
	tests := make([]testCaseG, 100)
	for i := range tests {
		n := rand.Intn(10) + 3
		tests[i] = testCaseG{n}
	}
	return tests
}

func validOutput(n int, nums []int64) bool {
	if len(nums) != n {
		return false
	}
	seen := make(map[int64]bool)
	xorOdd, xorEven := int64(0), int64(0)
	for i, v := range nums {
		if v < 0 || v >= 1<<31 {
			return false
		}
		if seen[v] {
			return false
		}
		seen[v] = true
		if (i+1)%2 == 1 {
			xorOdd ^= v
		} else {
			xorEven ^= v
		}
	}
	return xorOdd == xorEven
}

func parseOutput(out string) ([]int64, error) {
	fields := strings.Fields(out)
	nums := make([]int64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, err
		}
		nums[i] = v
	}
	return nums, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("1\n%d\n", tc.n)
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		nums, err := parseOutput(output)
		if err != nil {
			fmt.Printf("test %d: cannot parse output: %v\n", i+1, err)
			return
		}
		if !validOutput(tc.n, nums) {
			fmt.Printf("test %d failed:\ninput:%soutput:%s\n", i+1, input, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
