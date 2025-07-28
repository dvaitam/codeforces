package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n int
	k int
}

func possible(n int, k int) bool {
	for a := 0; a <= n; a++ {
		b := n - a
		cur := a*(a-1)/2 + b*(b-1)/2
		if cur == k {
			return true
		}
	}
	return false
}

func checkArray(arr []int, k int) bool {
	n := len(arr)
	pairs := 0
	for i := 0; i < n; i++ {
		if arr[i] != 1 && arr[i] != -1 {
			return false
		}
		for j := i + 1; j < n; j++ {
			if arr[i]*arr[j] == 1 {
				pairs++
			}
		}
	}
	return pairs == k
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(1)
	const cases = 100
	tests := make([]Test, cases)
	for i := range tests {
		n := rand.Intn(10) + 2
		maxK := n * (n - 1) / 2
		k := rand.Intn(maxK + 1)
		tests[i] = Test{n: n, k: k}
	}

	var input strings.Builder
	fmt.Fprintf(&input, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.k)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("error running binary:", err)
		fmt.Print(out.String())
		return
	}

	reader := bufio.NewReader(bytes.NewReader(out.Bytes()))
	for idx, tc := range tests {
		token, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("test %d: failed to read output\n", idx+1)
			return
		}
		token = strings.TrimSpace(token)
		if token != "YES" && token != "NO" {
			fmt.Printf("test %d: expected YES/NO got %s\n", idx+1, token)
			return
		}
		poss := possible(tc.n, tc.k)
		if token == "NO" {
			if poss {
				fmt.Printf("test %d: expected YES but got NO\n", idx+1)
				return
			}
			continue
		}
		// read n integers
		arr := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			if _, err := fmt.Fscan(reader, &arr[i]); err != nil {
				fmt.Printf("test %d: failed to read array\n", idx+1)
				return
			}
		}
		// consume remainder of line
		reader.ReadString('\n')
		if !poss {
			fmt.Printf("test %d: expected NO but got YES\n", idx+1)
			return
		}
		if !checkArray(arr, tc.k) {
			fmt.Printf("test %d: wrong array\n", idx+1)
			return
		}
	}
	fmt.Printf("verified %d test cases\n", len(tests))
}
