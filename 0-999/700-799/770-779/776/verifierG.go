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

func countDecrease(l, r uint64) uint64 {
	var cnt uint64
	for x := l; x <= r; x++ {
		y := x
		digits := make(map[uint64]struct{})
		for t := x; ; t >>= 4 {
			digits[t&0xF] = struct{}{}
			if t < 16 {
				break
			}
		}
		var sum uint64
		for d := range digits {
			sum += 1 << d
		}
		if (y ^ sum) < y {
			cnt++
		}
	}
	return cnt
}

func run(bin, in string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 1; i <= 100; i++ {
		l := uint64(i - 1)
		r := l + 15
		input := fmt.Sprintf("1\n%x %x\n", l, r)
		expected := fmt.Sprintf("%d", countDecrease(l, r))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d: execution error: %v\nOutput:\n%s\n", i, err, got)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("Test %d failed.\nInput:%sExpected:%s Got:%s\n", i, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
