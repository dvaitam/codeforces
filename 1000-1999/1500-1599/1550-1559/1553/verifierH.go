package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("time limit")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out)
	}
	return strings.TrimSpace(string(out)), nil
}

const inf int = 1 << 60

func solve(nums []int, b int) ([]int, []int, []int) {
	size := 1 << b
	ans := make([]int, size)
	mn := make([]int, size)
	mx := make([]int, size)
	if b == 0 {
		ans[0] = inf
		if len(nums) > 0 {
			mn[0] = 0
			mx[0] = 0
		} else {
			mn[0] = inf
			mx[0] = -inf
		}
		return ans, mn, mx
	}
	mid := 1 << (b - 1)
	leftNums := make([]int, 0)
	rightNums := make([]int, 0)
	for _, v := range nums {
		if v&mid == 0 {
			leftNums = append(leftNums, v)
		} else {
			rightNums = append(rightNums, v-mid)
		}
	}
	ansL, mnL, mxL := solve(leftNums, b-1)
	ansR, mnR, mxR := solve(rightNums, b-1)
	for mask := 0; mask < size; mask++ {
		hi := mask >> (b - 1)
		low := mask & (mid - 1)
		if hi == 0 {
			cross := inf
			if len(leftNums) > 0 && len(rightNums) > 0 {
				cross = (mnR[low] + mid) - mxL[low]
			}
			a := ansL[low]
			if ansR[low] < a {
				a = ansR[low]
			}
			if cross < a {
				a = cross
			}
			ans[mask] = a
			minVal := inf
			if len(leftNums) > 0 && mnL[low] < minVal {
				minVal = mnL[low]
			}
			if len(rightNums) > 0 && mnR[low]+mid < minVal {
				minVal = mnR[low] + mid
			}
			mn[mask] = minVal
			maxVal := -inf
			if len(leftNums) > 0 && mxL[low] > maxVal {
				maxVal = mxL[low]
			}
			if len(rightNums) > 0 && mxR[low]+mid > maxVal {
				maxVal = mxR[low] + mid
			}
			mx[mask] = maxVal
		} else {
			cross := inf
			if len(leftNums) > 0 && len(rightNums) > 0 {
				cross = (mnL[low] + mid) - mxR[low]
			}
			a := ansL[low]
			if ansR[low] < a {
				a = ansR[low]
			}
			if cross < a {
				a = cross
			}
			ans[mask] = a
			minVal := inf
			if len(rightNums) > 0 && mnR[low] < minVal {
				minVal = mnR[low]
			}
			if len(leftNums) > 0 && mnL[low]+mid < minVal {
				minVal = mnL[low] + mid
			}
			mn[mask] = minVal
			maxVal := -inf
			if len(rightNums) > 0 && mxR[low] > maxVal {
				maxVal = mxR[low]
			}
			if len(leftNums) > 0 && mxL[low]+mid > maxVal {
				maxVal = mxL[low] + mid
			}
			mx[mask] = maxVal
		}
	}
	return ans, mn, mx
}

func compute(input string) string {
	rdr := strings.NewReader(strings.TrimSpace(input) + "\n")
	var n, k int
	fmt.Fscan(rdr, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(rdr, &nums[i])
	}
	ans, _, _ := solve(nums, k)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	return sb.String()
}

func generateCases() []testCase {
	rand.Seed(8)
	cases := []testCase{}
	fixed := []struct {
		n, k int
		arr  []int
	}{
		{1, 1, []int{0}},
		{2, 2, []int{1, 2}},
	}
	for _, f := range fixed {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", f.n, f.k)
		for i, v := range f.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	for len(cases) < 100 {
		n := rand.Intn(4) + 1
		k := rand.Intn(4) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for i := 0; i < n; i++ {
			v := rand.Intn(1 << k)
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		inp := sb.String()
		cases = append(cases, testCase{inp, compute(inp)})
	}
	return cases
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierH.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCases()
	for i, tc := range cases {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != tc.expected {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%sexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
