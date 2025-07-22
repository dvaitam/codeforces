package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const mod = 1000000009

func fibSeq(n int) []int64 {
	f := make([]int64, n+2)
	f[1] = 1
	f[2] = 1
	for i := 3; i < len(f); i++ {
		f[i] = (f[i-1] + f[i-2]) % mod
	}
	return f
}

func generateCase(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	arr := make([]int64, n)
	for i := range arr {
		arr[i] = int64(rng.Intn(5))
	}
	fib := fibSeq(n + m + 5)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	outputs := make([]int64, 0, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("1 %d %d\n", l, r))
			for j := l - 1; j <= r-1; j++ {
				arr[j] = (arr[j] + fib[j-l+2]) % mod
			}
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
			sum := int64(0)
			for j := l - 1; j <= r-1; j++ {
				sum += arr[j]
				if sum >= mod {
					sum -= mod
				}
			}
			outputs = append(outputs, sum%mod)
		}
	}
	return sb.String(), outputs
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		lines := strings.Fields(out)
		if len(lines) != len(exp) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d lines got %d\ninput:\n%s", i+1, len(exp), len(lines), in)
			os.Exit(1)
		}
		for j, s := range lines {
			var got int64
			if _, err := fmt.Sscan(s, &got); err != nil || got != exp[j] {
				fmt.Fprintf(os.Stderr, "case %d: wrong answer on query %d expected %d got %s\ninput:\n%s", i+1, j+1, exp[j], s, in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
