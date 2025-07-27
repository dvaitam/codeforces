package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const mod int64 = 998244353

func expected(n int) []int64 {
	pow10 := make([]int64, n+1)
	pow10[0] = 1
	for i := 1; i <= n; i++ {
		pow10[i] = pow10[i-1] * 10 % mod
	}
	res := make([]int64, n)
	for i := 1; i <= n; i++ {
		if i == n {
			res[i-1] = 10
		} else {
			part1 := int64(2*10*9) * pow10[n-i-1] % mod
			ans := part1
			if n-i-1 > 0 {
				part2 := int64(n-i-1) * 10 % mod
				part2 = part2 * 9 % mod
				part2 = part2 * 9 % mod
				part2 = part2 * pow10[n-i-2] % mod
				ans = (ans + part2) % mod
			}
			res[i-1] = ans
		}
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(5)
	t := 100
	for idx := 0; idx < t; idx++ {
		n := rand.Intn(50) + 1
		wantArr := expected(n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("Runtime error on test %d: %v\n%s", idx+1, err, out.String())
			os.Exit(1)
		}
		gotLines := strings.Fields(strings.TrimSpace(out.String()))
		if len(gotLines) != n {
			fmt.Printf("Wrong answer on test %d: expected %d numbers got %d\n", idx+1, n, len(gotLines))
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			if gotLines[i] != fmt.Sprint(wantArr[i]) {
				fmt.Printf("Wrong answer on test %d position %d: expected %d got %s\n", idx+1, i+1, wantArr[i], gotLines[i])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed.")
}
