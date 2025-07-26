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

const mod int = 1000000007

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var s string
	fmt.Fscan(in, &s)
	cur := 0
	ans := 0
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == 'b' {
			cur++
			if cur >= mod {
				cur -= mod
			}
		} else {
			ans += cur
			if ans >= mod {
				ans -= mod
			}
			cur <<= 1
			if cur >= mod {
				cur %= mod
			}
		}
	}
	return fmt.Sprintf("%d", ans)
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
	r := rand.New(rand.NewSource(2))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		length := r.Intn(20) + 1
		var sb strings.Builder
		for j := 0; j < length; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('a')
			} else {
				sb.WriteByte('b')
			}
		}
		tests[i] = sb.String() + "\n"
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveB(t)
		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed. input: %sexpected %s got %s\n", i+1, t, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
