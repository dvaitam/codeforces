package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000009

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = (res * a) % mod
		}
		a = (a * a) % mod
		e >>= 1
	}
	return res
}

func solveCase(n int, m int64) string {
	if m < 63 {
		if int64(n)+1 > (int64(1) << m) {
			return "0"
		}
	}
	total := modPow(2, m)
	ans := int64(1)
	for i := 0; i < n; i++ {
		term := (total - 1 - int64(i)) % mod
		if term < 0 {
			term += mod
		}
		ans = (ans * term) % mod
	}
	return fmt.Sprintf("%d", ans)
}

type test struct{ input, expected string }

func generateTests() []test {
	rng := rand.New(rand.NewSource(42))
	var tests []test
	fixed := [][2]int{{1, 1}, {2, 1}, {3, 4}, {5, 6}}
	for _, f := range fixed {
		inp := fmt.Sprintf("%d %d\n", f[0], f[1])
		tests = append(tests, test{inp, solveCase(f[0], int64(f[1]))})
	}
	for len(tests) < 100 {
		n := rng.Intn(20) + 1
		m := rng.Intn(60) + 1
		inp := fmt.Sprintf("%d %d\n", n, m)
		tests = append(tests, test{inp, solveCase(n, int64(m))})
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
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
