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

func generateTests() [][]int {
	r := rand.New(rand.NewSource(44))
	tests := make([][]int, 100)
	for i := 0; i < 100; i++ {
		n := r.Intn(50) + 2
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(n) + 1
		}
		// ensure at least one duplicate
		arr[0] = arr[1]
		tests[i] = arr
	}
	return tests
}

func expected(arr []int) int {
	n := len(arr)
	freq := make(map[int]int)
	for _, v := range arr {
		freq[v]++
	}
	mX := 0
	for _, f := range freq {
		if f > mX {
			mX = f
		}
	}
	c := 0
	for _, f := range freq {
		if f == mX {
			c++
		}
	}
	ans := (n - c) / (mX - 1)
	ans--
	if ans < 0 {
		ans = 0
	}
	return ans
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	tests := generateTests()
	for i, arr := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(strconv.Itoa(len(arr)))
		sb.WriteByte('\n')
		for j, v := range arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		exp := expected(arr)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		var ans int
		if _, err := fmt.Sscan(got, &ans); err != nil {
			fmt.Printf("Test %d: cannot parse output %q\n", i+1, got)
			os.Exit(1)
		}
		if ans != exp {
			fmt.Printf("Test %d failed. Input:\n%sExpected %d got %d\n", i+1, sb.String(), exp, ans)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
