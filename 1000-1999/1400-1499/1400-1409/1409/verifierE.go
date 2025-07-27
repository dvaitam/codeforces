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

func solve(xs []int64, k int64) int {
	n := len(xs)
	sort.Slice(xs, func(i, j int) bool { return xs[i] < xs[j] })
	cnt := make([]int, n)
	r := 0
	for i := 0; i < n; i++ {
		for r < n && xs[r] <= xs[i]+k {
			r++
		}
		cnt[i] = r - i
	}
	suff := make([]int, n+1)
	for i := n - 1; i >= 0; i-- {
		if cnt[i] > suff[i+1] {
			suff[i] = cnt[i]
		} else {
			suff[i] = suff[i+1]
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		j := i + cnt[i]
		if j > n {
			j = n
		}
		total := cnt[i] + suff[j]
		if total > ans {
			ans = total
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	k := int64(rng.Intn(50) + 1)
	xs := make([]int64, n)
	ys := make([]int64, n)
	for i := 0; i < n; i++ {
		xs[i] = int64(rng.Intn(100)) + 1
		ys[i] = int64(rng.Intn(100)) + 1
	}
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	for i, v := range xs {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	for i, v := range ys {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	ans := solve(append([]int64(nil), xs...), k)
	expected := fmt.Sprintf("%d", ans)
	return input, expected
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %q got %q", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
