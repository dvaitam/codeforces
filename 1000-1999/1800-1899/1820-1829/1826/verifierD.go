package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solve(b []int64) int64 {
	n := len(b)
	const negInf int64 = -1 << 60
	pre := make([]int64, n+1)
	pre[0] = negInf
	for i := 1; i <= n; i++ {
		val := b[i-1] + int64(i)
		if val > pre[i-1] {
			pre[i] = val
		} else {
			pre[i] = pre[i-1]
		}
	}
	suf := make([]int64, n+2)
	suf[n+1] = negInf
	for i := n; i >= 1; i-- {
		val := b[i-1] - int64(i)
		if val > suf[i+1] {
			suf[i] = val
		} else {
			suf[i] = suf[i+1]
		}
	}
	ans := negInf
	for j := 2; j <= n-1; j++ {
		val := pre[j-1] + b[j-1] + suf[j+1]
		if val > ans {
			ans = val
		}
	}
	return ans
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

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierD <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(0)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 3
		b := make([]int64, n)
		for i := range b {
			b[i] = rand.Int63n(20) + 1
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i, v := range b {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		expected := fmt.Sprint(solve(b))
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s\n", t, err, output)
			os.Exit(1)
		}
		if output != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", t, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
