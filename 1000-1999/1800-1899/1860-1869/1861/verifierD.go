package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct{ input, expected string }

func lisLength(arr []int) int {
	dp := make([]int, 0)
	for _, v := range arr {
		l, r := 0, len(dp)
		for l < r {
			m := (l + r) / 2
			if dp[m] < v {
				l = m + 1
			} else {
				r = m
			}
		}
		if l == len(dp) {
			dp = append(dp, v)
		} else {
			dp[l] = v
		}
	}
	return len(dp)
}

func solveCase(arr []int) string {
	l := lisLength(arr)
	return fmt.Sprintf("%d", len(arr)-l)
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(45))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(10) + 1
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rng.Intn(20)
		}
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		tests = append(tests, test{sb.String(), solveCase(arr)})
	}
	return tests
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		got, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Printf("Runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(t.expected) {
			fmt.Printf("Wrong answer on test %d\nInput:%sExpected:%s\nGot:%s\n", i+1, t.input, t.expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
