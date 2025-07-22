package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func expected(n int, l, r, Ql, Qr int64, w []int64) int64 {
	pre := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pre[i+1] = pre[i] + w[i]
	}
	total := pre[n]
	ans := int64(1<<63 - 1)
	for L := 0; L <= n; L++ {
		leftCost := pre[L] * l
		rightCost := (total - pre[L]) * r
		R := n - L
		var penalty int64
		if L > R {
			diff := int64(L - R - 1)
			if diff > 0 {
				penalty = diff * Ql
			}
		} else if R > L {
			diff := int64(R - L - 1)
			if diff > 0 {
				penalty = diff * Qr
			}
		}
		cost := leftCost + rightCost + penalty
		if cost < ans {
			ans = cost
		}
	}
	return ans
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		l := rng.Int63n(100) + 1
		r := rng.Int63n(100) + 1
		Ql := rng.Int63n(20) + 1
		Qr := rng.Int63n(20) + 1
		w := make([]int64, n)
		for j := 0; j < n; j++ {
			w[j] = rng.Int63n(100) + 1
		}

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d %d %d\n", n, l, r, Ql, Qr))
		for j, val := range w {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(val, 10))
		}
		input.WriteByte('\n')

		expectedOut := strconv.FormatInt(expected(n, l, r, Ql, Qr, w), 10)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input.String())
			os.Exit(1)
		}
		if got != expectedOut {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expectedOut, got, input.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
