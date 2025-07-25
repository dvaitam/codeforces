package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testD struct {
	s string
}

func check(sum []int) (bool, string) {
	n := len(sum)
	ans := make([]byte, n)
	i := 0
	for i < n/2 {
		if sum[i] == sum[n-1-i] {
			i++
		} else if sum[i] == sum[n-1-i]+1 || sum[i] == sum[n-1-i]+11 {
			sum[i]--
			sum[i+1] += 10
		} else if sum[i] == sum[n-1-i]+10 {
			sum[n-2-i]--
			sum[n-1-i] += 10
		} else {
			return false, ""
		}
	}
	if n%2 == 1 {
		mid := n / 2
		if sum[mid]%2 != 0 || sum[mid] > 18 || sum[mid] < 0 {
			return false, ""
		}
		ans[mid] = byte(sum[mid]/2) + '0'
	}
	for j := 0; j < n/2; j++ {
		if sum[j] > 18 || sum[j] < 0 {
			return false, ""
		}
		ans[j] = byte((sum[j]+1)/2) + '0'
		ans[n-1-j] = byte(sum[j]/2) + '0'
	}
	if ans[0] <= '0' {
		return false, ""
	}
	return true, string(ans)
}

func solveD(s string) string {
	n := len(s)
	sum := make([]int, n)
	for i := 0; i < n; i++ {
		sum[i] = int(s[i] - '0')
	}
	if ok, res := check(append([]int{}, sum...)); ok {
		return res
	}
	if n > 1 && s[0] == '1' {
		n1 := n - 1
		sum1 := make([]int, n1)
		for i := 0; i < n1; i++ {
			sum1[i] = int(s[i+1] - '0')
		}
		sum1[0] += 10
		if ok, res := check(sum1); ok {
			return res
		}
	}
	return "0"
}

func genTests() []testD {
	rand.Seed(4)
	tests := make([]testD, 100)
	for i := range tests {
		l := rand.Intn(20) + 1
		b := make([]byte, l)
		for j := range b {
			if j == 0 {
				b[j] = byte(rand.Intn(9)+1) + '0'
			} else {
				b[j] = byte(rand.Intn(10)) + '0'
			}
		}
		tests[i] = testD{s: string(b)}
	}
	tests = append(tests, testD{s: "4"})
	tests = append(tests, testD{s: "11"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := fmt.Sprintf("%s\n", t.s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		expected := solveD(t.s)
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("Test %d failed\nInput:%sExpected %s Got %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
