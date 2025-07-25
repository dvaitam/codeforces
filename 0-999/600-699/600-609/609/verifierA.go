package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expected(n, m int, drives []int) int {
	sort.Slice(drives, func(i, j int) bool { return drives[i] > drives[j] })
	sum := 0
	cnt := 0
	for _, v := range drives {
		sum += v
		cnt++
		if sum >= m {
			break
		}
	}
	return cnt
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts)-2 < n {
			fmt.Printf("test %d invalid length\n", idx)
			os.Exit(1)
		}
		drives := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			drives[i] = v
		}
		exp := expected(n, m, drives)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d", n, m)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, " %d", drives[i])
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprint(exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
