package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod = 1000000007
const inv9 = 111111112

type test struct {
	input    string
	expected string
}

func solve(input string) string {
	s := strings.TrimSpace(input)
	n := len(s)
	digits := make([]int, n+1)
	for i := 1; i <= n; i++ {
		digits[i] = int(s[i-1] - '0')
	}
	pow10 := make([]int, n+2)
	pow10[0] = 1
	for i := 1; i <= n+1; i++ {
		pow10[i] = int(int64(pow10[i-1]) * 10 % mod)
	}
	P := make([]int, n+1)
	for i := 1; i <= n; i++ {
		P[i] = int((int64(P[i-1])*10 + int64(digits[i])) % mod)
	}
	Suf := make([]int, n+2)
	for i := n; i >= 1; i-- {
		Suf[i] = int((int64(digits[i])*int64(pow10[n-i]) + int64(Suf[i+1])) % mod)
	}
	var sum1, sum2 int64
	for l := 1; l <= n; l++ {
		x := (int64(pow10[n-l+1]) - 1 + mod) % mod
		x = x * inv9 % mod
		sum1 = (sum1 + int64(P[l-1])*x) % mod
	}
	for r := 1; r <= n; r++ {
		sum2 = (sum2 + int64(r)*int64(Suf[r+1])) % mod
	}
	res := (sum1 + sum2) % mod
	return fmt.Sprint(res)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(44))
	tests := []test{}
	fixed := []string{"0", "12345", "999"}
	for _, f := range fixed {
		tests = append(tests, test{f, solve(f)})
	}
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		var sb strings.Builder
		for i := 0; i < n; i++ {
			sb.WriteByte(byte('0' + rng.Intn(10)))
		}
		inp := sb.String()
		tests = append(tests, test{inp, solve(inp)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input + "\n")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%s\nExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
