package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		a = -a
	}
	return a
}

func checkCase(n int, arr []int64, out string) error {
	tokens := strings.Fields(out)
	if len(tokens) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("invalid k")
	}
	if k < 0 || k > n {
		return fmt.Errorf("k out of range")
	}
	if len(tokens) != 1+4*k {
		return fmt.Errorf("expected %d operation lines, got %d tokens", k, len(tokens))
	}
	idx := 1
	for op := 0; op < k; op++ {
		if idx+3 >= len(tokens) {
			return fmt.Errorf("not enough tokens for op %d", op+1)
		}
		i, _ := strconv.Atoi(tokens[idx])
		idx++
		j, _ := strconv.Atoi(tokens[idx])
		idx++
		x64, _ := strconv.ParseInt(tokens[idx], 10, 64)
		idx++
		y64, _ := strconv.ParseInt(tokens[idx], 10, 64)
		idx++
		if i < 1 || i > n || j < 1 || j > n || i == j {
			return fmt.Errorf("bad indices in op %d", op+1)
		}
		if x64 < 1 || x64 > 2_000_000_000 || y64 < 1 || y64 > 2_000_000_000 {
			return fmt.Errorf("values out of range in op %d", op+1)
		}
		m0 := arr[i-1]
		if arr[j-1] < m0 {
			m0 = arr[j-1]
		}
		m1 := x64
		if y64 < m1 {
			m1 = y64
		}
		if m0 != m1 {
			return fmt.Errorf("min mismatch in op %d", op+1)
		}
		arr[i-1] = x64
		arr[j-1] = y64
	}
	for i := 1; i < n; i++ {
		if gcd(arr[i-1], arr[i]) != 1 {
			return fmt.Errorf("array not good after ops")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		arr := make([]int64, n)
		for i := range arr {
			arr[i] = rand.Int63n(100) + 1
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range arr {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\n", t, err)
			os.Exit(1)
		}
		arrCopy := make([]int64, len(arr))
		copy(arrCopy, arr)
		if err := checkCase(n, arrCopy, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:%soutput:%s\n", t, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
