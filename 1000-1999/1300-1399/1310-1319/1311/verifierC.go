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

func expected(n, m int, s string, pos []int) string {
	cnt := make([]int, n+2)
	for _, p := range pos {
		if p <= n {
			cnt[p]++
		}
	}
	for i := n; i >= 1; i-- {
		cnt[i] += cnt[i+1]
	}
	ans := make([]int64, 26)
	for i := 1; i <= n; i++ {
		times := int64(cnt[i] + 1)
		ans[s[i-1]-'a'] += times
	}
	parts := make([]string, 26)
	for i, v := range ans {
		parts[i] = strconv.FormatInt(v, 10)
	}
	return strings.Join(parts, " ")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
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
		if len(fields) < 3 {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		s := fields[2]
		if len(s) != n {
			fmt.Printf("test %d malformed string length\n", idx)
			os.Exit(1)
		}
		if len(fields) != 3+m {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		pos := make([]int, m)
		for i := 0; i < m; i++ {
			v, _ := strconv.Atoi(fields[3+i])
			pos[i] = v
		}
		exp := expected(n, m, s, pos)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		input.WriteString(fmt.Sprintf("%s\n", s))
		for i := 0; i < m; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(pos[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
