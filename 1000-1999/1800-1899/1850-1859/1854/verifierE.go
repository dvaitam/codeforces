package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solveE(m int64) []int {
	cnt := make([]int64, 60)
	cnt[0] = 1
	ans := make([]int, 0)
	for i, r, c := 1, int64(0), 0; i <= 60 && c <= 60; {
		if r < (int64(1)<<uint(i)) && m >= cnt[60-i] {
			m -= cnt[60-i]
			for j := 59; j >= i; j-- {
				cnt[j] += cnt[j-i]
			}
			r++
			c++
			ans = append(ans, i)
		} else {
			r = 0
			i++
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func parse(out string) ([]int, error) {
	var k int
	r := strings.NewReader(out)
	if _, err := fmt.Fscan(r, &k); err != nil {
		return nil, err
	}
	res := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(r, &res[i]); err != nil {
			return nil, err
		}
	}
	return res, nil
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(4))
	const tests = 100
	for t := 0; t < tests; t++ {
		m := r.Int63n(1<<20) + 1
		input := fmt.Sprintf("%d\n", m)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d error: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := parse(out)
		if err != nil {
			fmt.Printf("test %d invalid output\n", t+1)
			os.Exit(1)
		}
		want := solveE(m)
		if !equal(got, want) {
			fmt.Printf("test %d failed: expected %v got %v\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
