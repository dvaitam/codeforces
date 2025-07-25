package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(n int, k int64, arr []int64) string {
	limit := int64(1e15)
	powers := []int64{}
	p := int64(1)
	for {
		powers = append(powers, p)
		if k == 1 {
			break
		}
		if k == -1 {
			if p == 1 {
				p = -1
				continue
			}
			break
		}
		if abs(p) > limit/abs(k) {
			break
		}
		p *= k
	}
	freq := map[int64]int64{0: 1}
	var sum, result int64
	for _, v := range arr {
		sum += v
		for _, pw := range powers {
			result += freq[sum-pw]
		}
		freq[sum]++
	}
	return fmt.Sprintf("%d", result)
}

func run(binary, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ks := []int64{-2, -1, 1, 2, 3}
	for i := 1; i <= 100; i++ {
		n := i%5 + 3
		k := ks[i%len(ks)]
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			arr[j] = int64((i+j)%5 - 2)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		inp := sb.String()
		exp := expected(n, k, arr)
		got, err := run(bin, inp)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", i, err, got)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:%s Got:%s\n", i, inp, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
