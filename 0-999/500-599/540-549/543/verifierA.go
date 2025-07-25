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

func solveA(n, m, b, mod int, a []int) int {
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, b+1)
	}
	dp[0][0] = 1
	for _, ai := range a {
		for j := 1; j <= m; j++ {
			for k := ai; k <= b; k++ {
				dp[j][k] += dp[j-1][k-ai]
				if dp[j][k] >= mod {
					dp[j][k] -= mod
				}
			}
		}
	}
	result := 0
	for k := 0; k <= b; k++ {
		result += dp[m][k]
		if result >= mod {
			result -= mod
		}
	}
	return result
}

func genTestA() (string, int) {
	n := rand.Intn(4) + 1
	m := rand.Intn(5) + 1
	b := rand.Intn(5) + 1
	mod := 1000000007
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rand.Intn(b + 1)
	}
	input := fmt.Sprintf("%d %d %d %d\n", n, m, b, mod)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", a[i])
	}
	input += "\n"
	expected := solveA(n, m, b, mod, a)
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= 100; t++ {
		input, expected := genTestA()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %d\nGot: %s\n", t, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
