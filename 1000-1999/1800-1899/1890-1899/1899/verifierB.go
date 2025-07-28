package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func segmentDiff(prefix []int64, k int) int64 {
	n := len(prefix) - 1
	maxSum := int64(math.MinInt64)
	minSum := int64(math.MaxInt64)
	for i := k; i <= n; i += k {
		sum := prefix[i] - prefix[i-k]
		if sum > maxSum {
			maxSum = sum
		}
		if sum < minSum {
			minSum = sum
		}
	}
	return maxSum - minSum
}

func solve(a []int64) int64 {
	n := len(a)
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + a[i]
	}
	var ans int64
	for d := 1; d*d <= n; d++ {
		if n%d == 0 {
			diff := segmentDiff(prefix, d)
			if diff > ans {
				ans = diff
			}
			if d != n/d {
				diff2 := segmentDiff(prefix, n/d)
				if diff2 > ans {
					ans = diff2
				}
			}
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		a := make([]int64, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Int63n(1000) + 1
			a[j] = val
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := fmt.Sprintf("%d\n", solve(a))
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%sq\n got:%q\n", i+1, input, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
