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

type testCase struct {
	nums   []int
	expect string
}

func solve(nums []int) string {
	for _, x := range nums {
		if x%2 != 0 {
			return "First"
		}
	}
	return "Second"
}

func generateTests() []testCase {
	r := rand.New(rand.NewSource(43))
	tests := make([]testCase, 0, 100)

	// fixed edge cases
	tests = append(tests, testCase{[]int{1}, solve([]int{1})})
	tests = append(tests, testCase{[]int{2}, solve([]int{2})})
	tests = append(tests, testCase{[]int{2, 4}, solve([]int{2, 4})})
	tests = append(tests, testCase{[]int{1, 2, 3}, solve([]int{1, 2, 3})})

	for len(tests) < 100 {
		n := r.Intn(20) + 1
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i] = r.Intn(1000)
		}
		tests = append(tests, testCase{nums, solve(nums)})
	}
	return tests
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Env = os.Environ()
	timer := time.AfterFunc(2*time.Second, func() { cmd.Process.Kill() })
	err := cmd.Run()
	timer.Stop()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	path := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d\n", len(tc.nums))
		for j, v := range tc.nums {
			if j+1 == len(tc.nums) {
				input += fmt.Sprintf("%d\n", v)
			} else {
				input += fmt.Sprintf("%d ", v)
			}
		}
		out, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nOutput: %s\n", i+1, err, out)
			os.Exit(1)
		}
		res := strings.TrimSpace(out)
		if res != tc.expect {
			fmt.Printf("Test %d failed\nInput: %sExpected: %s\nGot: %s\n", i+1, input, tc.expect, res)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
