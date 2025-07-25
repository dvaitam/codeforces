package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// solveA implements the reference solution for problem A.
func solveA(n, k int, s string) string {
	posG, posT := -1, -1
	for i, ch := range s {
		switch ch {
		case 'G':
			posG = i
		case 'T':
			posT = i
		}
	}
	if (posT-posG)%k != 0 {
		return "NO"
	}
	step := k
	if posT < posG {
		step = -k
	}
	for cur := posG; ; cur += step {
		if s[cur] == '#' {
			return "NO"
		}
		if cur == posT {
			return "YES"
		}
	}
}

// genTestA generates a random test case.
func genTestA() (int, int, string) {
	n := rand.Intn(99) + 2  // 2..100
	k := rand.Intn(n-1) + 1 // 1..n-1
	arr := make([]byte, n)
	for i := 0; i < n; i++ {
		if rand.Intn(3) == 0 {
			arr[i] = '#'
		} else {
			arr[i] = '.'
		}
	}
	g := rand.Intn(n)
	t := rand.Intn(n)
	for t == g {
		t = rand.Intn(n)
	}
	arr[g] = 'G'
	arr[t] = 'T'
	return n, k, string(arr)
}

func runBinary(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rand.Seed(time.Now().UnixNano())
	path := os.Args[1]
	for i := 0; i < 100; i++ {
		n, k, s := genTestA()
		input := fmt.Sprintf("%d %d\n%s\n", n, k, s)
		expected := solveA(n, k, s)
		got, err := runBinary(path, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
