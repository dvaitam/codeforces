package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(s string) int64 {
	n := len(s)
	f := make([]int, n)
	var ans int64
	const P uint64 = 131
	var h1, h2, pPow uint64 = 0, 0, 1
	for i := 0; i < n; i++ {
		c := uint64(s[i])
		h1 = h1*P + c
		h2 = h2 + c*pPow
		pPow *= P
		if h1 == h2 {
			halfIdx := ((i + 1) >> 1) - 1
			if halfIdx >= 0 {
				f[i] = f[halfIdx] + 1
			} else {
				f[i] = 1
			}
		}
		ans += int64(f[i])
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		exp := expected(line)
		var input bytes.Buffer
		input.WriteString(line)
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		var got int64
		outStr := strings.TrimSpace(string(out))
		fmt.Sscan(outStr, &got)
		if got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
