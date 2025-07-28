package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin string, input string) (string, error) {
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

func parseOps(out string, n int, limit int) ([][2]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return nil, fmt.Errorf("failed to read k: %v", err)
	}
	if k < 0 || k > limit {
		return nil, fmt.Errorf("k out of range: %d", k)
	}
	ops := make([][2]int, k)
	for i := 0; i < k; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return nil, fmt.Errorf("failed to read op %d: %v", i+1, err)
		}
		if x < 1 || x > n || y < 1 || y > n {
			return nil, fmt.Errorf("invalid indices in op %d", i+1)
		}
		ops[i] = [2]int{x - 1, y - 1}
	}
	return ops, nil
}

func applyOps(a []int, ops [][2]int) {
	for _, op := range ops {
		a[op[0]] += a[op[1]]
	}
}

func isNonDecreasing(a []int) bool {
	for i := 1; i < len(a); i++ {
		if a[i] < a[i-1] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	const tests = 100
	for t := 0; t < tests; t++ {
		n := 1 + r.Intn(20)
		a := make([]int, n)
		for i := range a {
			a[i] = r.Intn(41) - 20
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		out, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d error: %v\n", t+1, err)
			os.Exit(1)
		}
		ops, err := parseOps(out, n, 50)
		if err != nil {
			fmt.Printf("test %d invalid output: %v\n", t+1, err)
			os.Exit(1)
		}
		b := make([]int, len(a))
		copy(b, a)
		applyOps(b, ops)
		if !isNonDecreasing(b) {
			fmt.Printf("test %d failed: array not sorted after operations\n", t+1)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
