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

type TestA struct {
	l, r, x, y, k int64
}

const maxA = 10000000

func generateTests() []TestA {
	r := rand.New(rand.NewSource(1))
	tests := make([]TestA, 0, 120)
	for i := 0; i < 120; i++ {
		l := int64(r.Intn(maxA) + 1)
		rVal := int64(r.Intn(maxA) + 1)
		if l > rVal {
			l, rVal = rVal, l
		}
		x := int64(r.Intn(maxA) + 1)
		y := int64(r.Intn(maxA) + 1)
		if x > y {
			x, y = y, x
		}
		k := int64(r.Intn(maxA) + 1)
		tests = append(tests, TestA{l, rVal, x, y, k})
	}
	return tests
}

func expected(t TestA) string {
	low := (t.l + t.k - 1) / t.k
	high := t.r / t.k
	if low <= high && max64(low, t.x) <= min64(high, t.y) {
		return "YES"
	}
	return "NO"
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("timeout")
	}
	if err != nil {
		return "", fmt.Errorf("%v: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d %d %d\n", t.l, t.r, t.x, t.y, t.k)
		want := expected(t)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got = strings.TrimSpace(strings.ToUpper(got))
		if got != want {
			fmt.Printf("Test %d failed:\nInput:%sExpected: %s\nGot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed!\n", len(tests))
}
