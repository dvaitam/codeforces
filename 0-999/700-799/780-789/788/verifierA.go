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

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func solveA(nums []int64) int64 {
	n := len(nums)
	if n < 2 {
		return 0
	}
	diff := make([]int64, n-1)
	for i := 0; i < n-1; i++ {
		diff[i] = abs64(nums[i+1] - nums[i])
	}
	dpPlus := make([]int64, n-1)
	dpMinus := make([]int64, n-1)
	dpPlus[0] = diff[0]
	dpMinus[0] = -diff[0]
	ans := dpPlus[0]
	if dpMinus[0] > ans {
		ans = dpMinus[0]
	}
	for i := 1; i < n-1; i++ {
		if dpMinus[i-1]+diff[i] > diff[i] {
			dpPlus[i] = dpMinus[i-1] + diff[i]
		} else {
			dpPlus[i] = diff[i]
		}
		if dpPlus[i-1]-diff[i] > -diff[i] {
			dpMinus[i] = dpPlus[i-1] - diff[i]
		} else {
			dpMinus[i] = -diff[i]
		}
		if dpPlus[i] > ans {
			ans = dpPlus[i]
		}
		if dpMinus[i] > ans {
			ans = dpMinus[i]
		}
	}
	return ans
}

type testCase struct {
	nums []int64
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 2
	nums := make([]int64, n)
	for i := range nums {
		nums[i] = rng.Int63n(2000) - 1000
	}
	return testCase{nums}
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tc.nums))
	for i, v := range tc.nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func expected(tc testCase) string {
	return fmt.Sprintf("%d", solveA(tc.nums))
}

func runCase(bin string, tc testCase) error {
	input := formatInput(tc)
	exp := expected(tc)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	if out != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %s got %s\ninput:\n%s", exp, out, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{nums: []int64{1, 2}},
		{nums: []int64{-5, 0, 5}},
		{nums: []int64{5, 4, 3, 2, 1}},
		{nums: []int64{0, 0, 0}},
	}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
