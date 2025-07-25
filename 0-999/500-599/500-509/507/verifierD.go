package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type TestD struct {
	n int
	k int
	m int
}

func generateTests() []TestD {
	rand.Seed(4)
	tests := make([]TestD, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		k := rand.Intn(9) + 2
		m := rand.Intn(1000) + 1
		tests[i] = TestD{n, k, m}
	}
	return tests
}

func expected(t TestD) int {
	n, k, m := t.n, t.k, t.m
	pow10 := make([]int, n)
	pow10[0] = 1 % k
	for i := 1; i < n; i++ {
		pow10[i] = pow10[i-1] * 10 % k
	}
	dpCurr := make([][2]int, k)
	dpNext := make([][2]int, k)
	dpCurr[0][0] = 1 % m
	for i := 0; i < n; i++ {
		for r := 0; r < k; r++ {
			dpNext[r][0] = 0
			dpNext[r][1] = 0
		}
		start := 0
		if i+1 == n {
			start = 1
		}
		for r := 0; r < k; r++ {
			for flag := 0; flag < 2; flag++ {
				v := dpCurr[r][flag]
				if v == 0 {
					continue
				}
				for d := start; d <= 9; d++ {
					newR := (d*pow10[i] + r) % k
					newFlag := flag
					if flag == 0 && newR == 0 && d != 0 {
						newFlag = 1
					}
					dpNext[newR][newFlag] += v
					if dpNext[newR][newFlag] >= m {
						dpNext[newR][newFlag] %= m
					}
				}
			}
		}
		dpCurr, dpNext = dpNext, dpCurr
	}
	ans := 0
	for r := 0; r < k; r++ {
		ans += dpCurr[r][1]
		if ans >= m {
			ans %= m
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
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		input := fmt.Sprintf("%d %d %d\n", t.n, t.k, t.m)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		val, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Printf("test %d: invalid output\n", i+1)
			os.Exit(1)
		}
		exp := expected(t)
		if val != exp {
			fmt.Printf("test %d: expected %d got %d\n", i+1, exp, val)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
