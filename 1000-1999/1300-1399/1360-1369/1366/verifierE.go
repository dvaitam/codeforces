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

const mod int64 = 998244353

func solveCase(n, m int, a, b []int) int64 {
	idx := n - 1
	ans := int64(1)
	for i := m - 1; i >= 0; i-- {
		last := -1
		for idx >= 0 && a[idx] >= b[i] {
			if a[idx] == b[i] && last == -1 {
				last = idx
			}
			idx--
		}
		if last == -1 {
			return 0
		}
		if i == 0 {
			if idx != -1 {
				return 0
			}
		} else {
			ans = ans * int64(last-idx) % mod
		}
	}
	return ans % mod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
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
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("invalid test %d\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields) != 2+n+m {
			fmt.Printf("test %d wrong count\n", idx)
			os.Exit(1)
		}
		a := make([]int, n)
		b := make([]int, m)
		pos := 2
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[pos])
			a[i] = v
			pos++
		}
		for i := 0; i < m; i++ {
			v, _ := strconv.Atoi(fields[pos])
			b[i] = v
			pos++
		}
		expected := solveCase(n, m, a, b)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", a[i]))
		}
		input.WriteByte('\n')
		for i := 0; i < m; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", b[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
