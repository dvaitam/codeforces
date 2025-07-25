package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expected(n int, p []int64, s string) int64 {
	prefixA := make([]int64, n+1)
	prefixB := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefixA[i+1] = prefixA[i]
		prefixB[i+1] = prefixB[i]
		if s[i] == 'A' {
			prefixA[i+1] += p[i]
		} else {
			prefixB[i+1] += p[i]
		}
	}
	totalA := prefixA[n]
	totalB := prefixB[n]
	ans := totalB
	for k := 0; k <= n; k++ {
		val := totalB - prefixB[k] + prefixA[k]
		if val > ans {
			ans = val
		}
		val = prefixB[n-k] + (totalA - prefixA[n-k])
		if val > ans {
			ans = val
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		pArr := make([]int64, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			pArr[i] = int64(rand.Intn(1000) + 1)
			sb.WriteString(fmt.Sprintf("%d ", pArr[i]))
		}
		sb.WriteByte('\n')
		var s strings.Builder
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				s.WriteByte('A')
			} else {
				s.WriteByte('B')
			}
		}
		str := s.String()
		sb.WriteString(str)
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(n, pArr, str)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", t+1, err, input, out)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.Fields(out)[0], 10, 64)
		if err != nil || val != exp {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %d\nGot: %s\n", t+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
