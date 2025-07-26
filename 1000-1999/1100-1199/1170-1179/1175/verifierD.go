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

func solve(n, k int, a []int64) int64 {
	suff := make([]int64, n)
	var sum int64
	for i := n - 1; i >= 0; i-- {
		sum += a[i]
		suff[i] = sum
	}
	ans := suff[0]
	if k > 1 {
		vals := make([]int64, n-1)
		for i := 1; i < n; i++ {
			vals[i-1] = suff[i]
		}
		sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
		for i := 0; i < k-1; i++ {
			ans += vals[i]
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := rng.Intn(n) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(100)) - 50
	}
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			input.WriteByte(' ')
		}
		fmt.Fprintf(&input, "%d", arr[i])
	}
	input.WriteByte('\n')
	out := fmt.Sprintf("%d\n", solve(n, k, arr))
	return input.String(), out
}

func runCase(bin string, in, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(exp) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s", exp, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
