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

func solveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		c := []int{}
		for _, v := range b {
			if len(c) == 0 || c[len(c)-1] != v {
				c = append(c, v)
			}
		}
		inc, dec := true, true
		for i := 1; i < len(c); i++ {
			if c[i] < c[i-1] {
				inc = false
			}
			if c[i] > c[i-1] {
				dec = false
			}
		}
		if inc || dec {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}
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
	r := rand.New(rand.NewSource(4))
	tests := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", r.Intn(21)-10))
		}
		sb.WriteByte('\n')
		tests[i] = sb.String()
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		expected := solveD(t)
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
