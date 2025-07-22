package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MODD = 1000003

func solveD(s string) int {
	n := len(s)
	non := make([]int, n+1)
	for i := 0; i < n; i++ {
		non[i+1] = non[i]
		if s[i] < '0' || s[i] > '9' {
			non[i+1]++
		}
	}
	dp := make([][]int, n)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for length := 1; length <= n; length++ {
		for i := 0; i+length-1 < n; i++ {
			j := i + length - 1
			if non[j+1]-non[i] == 0 {
				dp[i][j] = 1
				continue
			}
			ways := 0
			if (s[i] == '+' || s[i] == '-') && i+1 <= j {
				ways = dp[i+1][j]
			}
			for k := i + 1; k < j; k++ {
				c := s[k]
				if c == '+' || c == '-' || c == '*' || c == '/' {
					left := dp[i][k-1]
					if left != 0 {
						right := dp[k+1][j]
						if right != 0 {
							ways = (ways + left*right) % MODD
						}
					}
				}
			}
			dp[i][j] = ways
		}
	}
	return dp[0][n-1] % MODD
}

func genCase() (string, int) {
	length := rand.Intn(10) + 1
	chars := []byte("0123456789+-*/")
	b := make([]byte, length)
	for i := 0; i < length; i++ {
		b[i] = chars[rand.Intn(len(chars))]
	}
	s := string(b)
	expect := solveD(s)
	return s + "\n", expect
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(42)
	for t := 1; t <= 100; t++ {
		in, expect := genCase()
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			fmt.Println(in)
			return
		}
		if strings.TrimSpace(got) != fmt.Sprint(expect) {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d\nGot: %s\n", t, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
