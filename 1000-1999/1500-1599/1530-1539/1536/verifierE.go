package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 1000000007

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		b >>= 1
	}
	return res
}

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		var hashCount int64
		zeroExists := false
		for i := 0; i < n; i++ {
			var row string
			fmt.Fscan(in, &row)
			for _, ch := range row {
				if ch == '#' {
					hashCount++
				} else if ch == '0' {
					zeroExists = true
				}
			}
		}
		ans := modPow(2, hashCount)
		if !zeroExists {
			ans--
			if ans < 0 {
				ans += mod
			}
		}
		out.WriteString(fmt.Sprintf("%d\n", ans%mod))
	}
	return strings.TrimSpace(out.String())
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []string {
	r := rand.New(rand.NewSource(5))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(5) + 1
		m := r.Intn(5) + 1
		if n*m < 2 {
			n = 1
			m = 2
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for x := 0; x < n; x++ {
			for y := 0; y < m; y++ {
				v := r.Intn(3)
				if v == 0 {
					sb.WriteByte('#')
				} else if v == 1 {
					sb.WriteByte('0')
				} else {
					sb.WriteByte('.')
				}
			}
			sb.WriteByte('\n')
		}
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveE(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed. input: %sexpected %s got %s\n", i+1, t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
