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

func solve(nums []int) int64 {
	maxVal := 0
	for _, v := range nums {
		if v > maxVal {
			maxVal = v
		}
	}
	freq := make([]int64, maxVal+1)
	for _, v := range nums {
		freq[v]++
	}
	dp := make([]int64, maxVal+2)
	if maxVal >= 1 {
		dp[1] = freq[1]
	}
	for i := 2; i <= maxVal; i++ {
		take := dp[i-2] + freq[i]*int64(i)
		if dp[i-1] > take {
			dp[i] = dp[i-1]
		} else {
			dp[i] = take
		}
	}
	return dp[maxVal]
}

func buildCase(nums []int) (string, int64) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(nums)))
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), solve(nums)
}

func generateCase(rng *rand.Rand) (string, int64) {
	n := rng.Intn(20) + 1
	nums := make([]int, n)
	for i := range nums {
		nums[i] = rng.Intn(20) + 1
	}
	return buildCase(nums)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []struct {
		input string
		want  int64
	}
	// simple predetermined cases
	cases = append(cases, func() struct {
		input string
		want  int64
	} { nums := []int{1}; in, w := buildCase(nums); return struct {
		input string
		want  int64
	}{in, w} }())
	cases = append(cases, func() struct {
		input string
		want  int64
	} { nums := []int{1, 2, 3}; in, w := buildCase(nums); return struct {
		input string
		want  int64
	}{in, w} }())
	for i := 0; i < 100; i++ {
		in, w := generateCase(rng)
		cases = append(cases, struct {
			input string
			want  int64
		}{in, w})
	}

	for idx, tc := range cases {
		gotStr, err := runCandidate(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: bad output %q\n", idx+1, gotStr)
			os.Exit(1)
		}
		if got != tc.want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, tc.want, got, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
