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

func isPrime(x int64) bool {
	if x < 2 {
		return false
	}
	for i := int64(2); i*i <= x; i++ {
		if x%i == 0 {
			return false
		}
	}
	return true
}

func solve(n int64) (int, []int64) {
	var a, b int64
	for i := n; i >= 2; i-- {
		if isPrime(i) {
			a = i
			break
		}
	}
	for i := int64(2); i <= n; i++ {
		if !isPrime(i) {
			continue
		}
		sum := a + i
		if sum != n && !isPrime(n-sum) {
			continue
		}
		b = i
		break
	}
	if a == n {
		return 1, []int64{a}
	} else if a+b == n {
		return 2, []int64{a, b}
	}
	return 3, []int64{a, b, n - (a + b)}
}

func run(bin, input string) (string, error) {
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(out string) (int, []int64, error) {
	reader := strings.NewReader(out)
	var k int
	if _, err := fmt.Fscan(reader, &k); err != nil {
		return 0, nil, fmt.Errorf("cannot parse k: %v", err)
	}
	res := make([]int64, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &res[i]); err != nil {
			return 0, nil, fmt.Errorf("cannot parse value: %v", err)
		}
	}
	return k, res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := make([]int64, 0, 100)
	for i := int64(3); i <= 11; i += 2 {
		cases = append(cases, i)
	}
	for len(cases) < 100 {
		n := rng.Int63n(99999) + 3
		if n%2 == 0 {
			n++
		}
		cases = append(cases, n)
	}

	for idx, n := range cases {
		input := fmt.Sprintf("%d\n", n)
		expK, expVals := solve(n)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", idx+1, err, input)
			os.Exit(1)
		}
		k, vals, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v\noutput:%s", idx+1, err, out)
			os.Exit(1)
		}
		if k != len(vals) {
			fmt.Fprintf(os.Stderr, "case %d failed: k=%d len(vals)=%d\n", idx+1, k, len(vals))
			os.Exit(1)
		}
		if k != expK {
			fmt.Fprintf(os.Stderr, "case %d failed: expected k=%d got %d\n", idx+1, expK, k)
			os.Exit(1)
		}
		for i := 0; i < k; i++ {
			if vals[i] != expVals[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", idx+1, expVals, vals)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
