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

type testCase struct {
	input  string
	expect int64
}

func expected(a []int64, s string) int64 {
	n := len(a)
	pre := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		pre[i] = pre[i-1] + a[i-1]
	}
	suf := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suf[i] = suf[i+1]
		if s[i] == '1' {
			suf[i] += a[i]
		}
	}
	ans := suf[0]
	for i := n - 1; i >= 0; i-- {
		if s[i] == '1' {
			cand := suf[i+1] + pre[i]
			if cand > ans {
				ans = cand
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(1000))
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", a[i])
	}
	sb.WriteByte('\n')
	bits := make([]byte, n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 1 {
			bits[i] = '1'
		} else {
			bits[i] = '0'
		}
	}
	s := string(bits)
	sb.WriteString(s)
	sb.WriteByte('\n')
	ans := expected(a, s)
	return testCase{input: sb.String(), expect: ans}
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
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
		tc := generateCase(rng)
		out, err := run(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		var val int64
		if _, err := fmt.Sscan(out, &val); err != nil || val != tc.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", i+1, tc.expect, out, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
