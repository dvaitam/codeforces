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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func solveC(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		var s string
		fmt.Fscan(in, &s)
		dCount := 0
		kCount := 0
		mp := make(map[[2]int]int)
		for i := 0; i < n; i++ {
			if s[i] == 'D' {
				dCount++
			} else {
				kCount++
			}
			g := gcd(dCount, kCount)
			key := [2]int{dCount / g, kCount / g}
			mp[key]++
			if i > 0 {
				out.WriteByte(' ')
			}
			out.WriteString(fmt.Sprintf("%d", mp[key]))
		}
		out.WriteByte('\n')
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
	r := rand.New(rand.NewSource(3))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if r.Intn(2) == 0 {
				sb.WriteByte('D')
			} else {
				sb.WriteByte('K')
			}
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveC(t)
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
