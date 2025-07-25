package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func brute(n, k int, costs []int) []int {
	maxA := 0
	for _, v := range costs {
		if v > maxA {
			maxA = v
		}
	}
	limit := maxA * k
	dp := make([]bool, limit+1)
	dp[0] = true
	for step := 0; step < k; step++ {
		next := make([]bool, limit+1)
		for _, c := range costs {
			for s := 0; s <= limit-c; s++ {
				if dp[s] {
					next[s+c] = true
				}
			}
		}
		dp = next
	}
	res := []int{}
	for i := 0; i <= limit; i++ {
		if dp[i] {
			res = append(res, i)
		}
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		n := rand.Intn(4) + 1
		k := rand.Intn(4) + 1
		costs := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			costs[i] = rand.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("%d ", costs[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := brute(n, k, costs)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s\nOutput:\n%s\n", t+1, err, input, out)
			os.Exit(1)
		}
		outFields := strings.Fields(out)
		got := []int{}
		for _, f := range outFields {
			v, err := strconv.Atoi(f)
			if err != nil {
				fmt.Printf("Test %d failed: invalid number\n", t+1)
				os.Exit(1)
			}
			got = append(got, v)
		}
		sort.Ints(got)
		if len(got) != len(exp) {
			fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %v\nGot: %v\n", t+1, input, exp, got)
			os.Exit(1)
		}
		for i := range exp {
			if got[i] != exp[i] {
				fmt.Printf("Test %d failed\nInput:\n%s\nExpected: %v\nGot: %v\n", t+1, input, exp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed!")
}
