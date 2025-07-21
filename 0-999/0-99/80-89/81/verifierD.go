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

func expected(n, m int, a []int) string {
	v := make([]int, m+1)
	copy(v[1:], a)
	ans := make([]int, n)
	last := -1
	for i := 0; i < n; i++ {
		id := 0
		for j := 1; j <= m; j++ {
			if j == last || (i == n-1 && j == ans[0]) {
				continue
			}
			if v[j] > v[id] || (v[j] == v[id] && j == ans[0]) {
				id = j
			}
		}
		if v[id] == 0 {
			return "-1"
		}
		v[id]--
		ans[i] = id
		last = id
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ans[i]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		mVal, _ := strconv.Atoi(parts[1])
		if len(parts) != mVal+2 {
			fmt.Printf("test %d: wrong number of album sizes\n", idx)
			os.Exit(1)
		}
		arr := make([]int, mVal)
		for i := 0; i < mVal; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			arr[i] = v
		}
		expect := expected(nVal, mVal, arr)
		input := fmt.Sprintf("%d %d\n%s\n", nVal, mVal, strings.Join(parts[2:], " "))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %s\n     got: %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
