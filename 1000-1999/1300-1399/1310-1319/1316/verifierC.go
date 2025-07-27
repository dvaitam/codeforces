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

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
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

func firstNotDivisible(arr []int64, p int64) int {
	for i, v := range arr {
		if v%p != 0 {
			return i
		}
	}
	return len(arr) - 1
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	primes := []int64{2, 3, 5, 7, 11, 13}
	p := primes[rng.Intn(len(primes))]
	fa := make([]int64, n)
	fb := make([]int64, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
	good := false
	for i := 0; i < n; i++ {
		fa[i] = int64(rng.Intn(20) + 1)
		if fa[i]%p != 0 {
			good = true
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(fa[i]))
	}
	if !good {
		fa[0]++
		sb.Reset()
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(fa[i]))
		}
	}
	sb.WriteByte('\n')
	good = false
	for i := 0; i < m; i++ {
		fb[i] = int64(rng.Intn(20) + 1)
		if fb[i]%p != 0 {
			good = true
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(fb[i]))
	}
	if !good {
		fb[0]++
		sb.Reset()
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(fa[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(fb[i]))
		}
	}
	sb.WriteByte('\n')
	expect := firstNotDivisible(fa, p) + firstNotDivisible(fb, p)
	input := "1\n" + sb.String()
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, expect := generateCase(rng)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: no output\ninput:%s", t+1, input)
			os.Exit(1)
		}
		val, err := strconv.Atoi(fields[0])
		if err != nil || val != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:%s", t+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
