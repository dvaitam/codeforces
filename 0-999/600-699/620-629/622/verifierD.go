package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveD(n int) []int {
	res := make([]int, 0, 2*n)
	x := (n - 1) % 2
	if x == 0 {
		x = 2
	}
	for x != n+1 {
		res = append(res, x)
		x += 2
	}
	x -= 2
	for x > 0 {
		res = append(res, x)
		x -= 2
	}
	res = append(res, n)
	x = n % 2
	if x == 0 {
		x = 2
	}
	for x != n+2 {
		res = append(res, x)
		x += 2
	}
	x -= 4
	for x > 0 {
		res = append(res, x)
		x -= 2
	}
	return res
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(4)
	var tests []int
	tests = append(tests, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	for len(tests) < 100 {
		tests = append(tests, rand.Intn(50)+1)
	}
	for i, n := range tests {
		input := fmt.Sprintf("%d\n", n)
		res := solveD(n)
		var sb strings.Builder
		for j, v := range res {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		expected := sb.String()
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("test %d failed: n=%d expected\n%s\ngot\n%s\n", i+1, n, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
