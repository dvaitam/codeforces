package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const mod = 1000000007

func expected(n int) int64 {
	var f0, f1 int64
	for i := 1; i <= n; i++ {
		add := int64(1)
		if i%2 == 0 {
			add = (add + f1) % mod
			f0 = (f0 + add) % mod
		} else {
			add = (add + f0) % mod
			f1 = (f1 + add) % mod
		}
	}
	return (f0 + f1) % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.Atoi(line)
		exp := expected(n)
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
