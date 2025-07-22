package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCase struct {
	n      int64
	expect int64
}

func genSuperLucky(L int, cur int64, c4, c7 int, res *[]int64) {
	if L == 0 {
		if c4 == c7 {
			*res = append(*res, cur)
		}
		return
	}
	genSuperLucky(L-1, cur*10+4, c4+1, c7, res)
	genSuperLucky(L-1, cur*10+7, c4, c7+1, res)
}

func nextSuperLucky(n int64) int64 {
	nStr := fmt.Sprintf("%d", n)
	for L := len(nStr); L <= 18; L++ {
		if L%2 == 1 {
			continue
		}
		var nums []int64
		genSuperLucky(L, 0, 0, 0, &nums)
		sort.Slice(nums, func(i, j int) bool { return nums[i] < nums[j] })
		for _, v := range nums {
			if int64(len(fmt.Sprintf("%d", v))) != int64(L) {
				continue
			}
			if L > len(nStr) || v >= n {
				return v
			}
		}
	}
	// fallback shouldn't happen for constraints
	return 0
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 100)
	edges := []int64{1, 4, 7, 44, 47, 74, 77, 4444, 4477, 4747, 7777, 1000000000}
	for _, n := range edges {
		tests = append(tests, testCase{n: n, expect: nextSuperLucky(n)})
	}
	for len(tests) < 100 {
		n := rng.Int63n(1_000_000_000) + 1
		tests = append(tests, testCase{n: n, expect: nextSuperLucky(n)})
	}
	return tests
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%d\n", tc.n))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Sscan(out.String(), &got); err != nil {
		return fmt.Errorf("invalid output: %v", err)
	}
	if got != tc.expect {
		return fmt.Errorf("expected %d got %d", tc.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %d\n", i+1, err, tc.n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
