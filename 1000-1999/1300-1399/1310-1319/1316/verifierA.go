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

func expected(n int, m int64, arr []int64) int64 {
	var sum int64
	for _, v := range arr {
		sum += v
	}
	if sum < m {
		return sum
	}
	return m
}

func generateCase(rng *rand.Rand) (string, int, int64, []int64) {
	n := rng.Intn(20) + 1
	m := rng.Int63n(100000) + 1
	arr := make([]int64, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(m + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), n, m, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, n, m, arr := generateCase(rng)
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
		val, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: cannot parse output\ninput:%soutput:%s\n", t+1, input, out)
			os.Exit(1)
		}
		exp := expected(n, m, arr)
		if val != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:%s", t+1, exp, val, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
