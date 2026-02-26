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

func solveB(a []int) int64 {
	n := len(a)
	sumA := 0
	for _, v := range a {
		sumA += v
	}

	dp := make([]bool, sumA+2)
	dp[1] = true

	for i := 1; i <= n; i++ {
		if a[i-1] == 0 {
			continue
		}
		for j := sumA + 1; j >= i; j-- {
			if dp[j] {
				if j+a[i-1] <= sumA+1 {
					dp[j+a[i-1]] = true
				}
			}
		}
	}

	P := make([]int, n+1)
	for i := 1; i <= n; i++ {
		P[i] = P[i-1] + a[i-1]
	}

	var maxVP int64 = -1
	for j := 1; j <= sumA+1; j++ {
		if dp[j] {
			var vp int
			if j <= n {
				vp = P[j] - j + 1
			} else {
				vp = P[n] - j + 1
			}
			if int64(vp) > maxVP {
				maxVP = int64(vp)
			}
		}
	}
	return maxVP
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(2))
	const tests = 100
	for t := 0; t < tests; t++ {
		n := 1 + r.Intn(50)
		a := make([]int, n)
		for i := range a {
			a[i] = r.Intn(n + 1)
		}
		var sb strings.Builder
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
		var got int64
		if _, err := fmt.Sscanf(out, "%d", &got); err != nil {
			fmt.Printf("test %d invalid output\n", t+1)
			os.Exit(1)
		}
		want := solveB(a)
		if got != want {
			fmt.Printf("test %d failed: expected %d got %d\n", t+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
