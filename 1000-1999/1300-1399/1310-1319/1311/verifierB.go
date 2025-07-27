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

func expected(n, m int, arr []int, pos []int) string {
	b := make([]bool, n-1)
	for _, p := range pos {
		if p >= 1 && p < n {
			b[p-1] = true
		}
	}
	i := 0
	for i < n-1 {
		if !b[i] {
			i++
			continue
		}
		l := i
		for i < n-1 && b[i] {
			i++
		}
		sort.Ints(arr[l : i+1])
	}
	ok := true
	for i := 0; i < n-1; i++ {
		if arr[i] > arr[i+1] {
			ok = false
			break
		}
	}
	if ok {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB.txt")
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
		if len(fields) < 2 {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		need := 2 + n + m
		if len(fields) != need {
			fmt.Printf("test %d malformed\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			arr[i] = v
		}
		pos := make([]int, m)
		for i := 0; i < m; i++ {
			v, _ := strconv.Atoi(fields[2+n+i])
			pos[i] = v
		}
		exp := expected(n, m, append([]int(nil), arr...), pos)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')
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
