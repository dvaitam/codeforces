package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

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
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n, k int, sizes []int) int {
	if n <= k {
		return sizes[n-1]
	}
	single := 2*k - n
	maxVal := 0
	for i := 0; i < single; i++ {
		if sizes[i] > maxVal {
			maxVal = sizes[i]
		}
	}
	left, right := single, n-1
	for left < right {
		sum := sizes[left] + sizes[right]
		if sum > maxVal {
			maxVal = sum
		}
		left++
		right--
	}
	return maxVal
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(100) + 1
	kMin := (n + 1) / 2
	k := rng.Intn(100-kMin+1) + kMin
	sizes := make([]int, n)
	for i := 0; i < n; i++ {
		sizes[i] = rng.Intn(1000000) + 1
	}
	sort.Ints(sizes)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", sizes[i])
	}
	sb.WriteByte('\n')
	exp := fmt.Sprintf("%d", expected(n, k, sizes))
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	fixed := []struct {
		n, k  int
		sizes []int
	}{
		{2, 1, []int{1, 2}},
		{3, 2, []int{2, 3, 5}},
		{4, 2, []int{3, 5, 7, 9}},
	}

	caseNum := 1
	for _, tc := range fixed {
		sizes := append([]int(nil), tc.sizes...)
		sort.Ints(sizes)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", sizes[i])
		}
		sb.WriteByte('\n')
		exp := fmt.Sprintf("%d", expected(tc.n, tc.k, sizes))
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum, err, sb.String())
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum, exp, out, sb.String())
			os.Exit(1)
		}
		caseNum++
	}

	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", caseNum, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", caseNum, exp, out, in)
			os.Exit(1)
		}
		caseNum++
	}
	fmt.Println("All tests passed")
}
