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

func expected(b []int, queries [][2]int) []int {
	ans := make([]int, len(queries))
	for i, q := range queries {
		set := make(map[int]struct{})
		for j := q[0] - 1; j <= q[1]-1; j++ {
			set[b[j]] = struct{}{}
		}
		ans[i] = len(set)
	}
	return ans
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
		m, _ := strconv.Atoi(parts[0])
		if len(parts) < 1+m+1 {
			fmt.Printf("test %d: short line\n", idx)
			os.Exit(1)
		}
		arr := make([]int, m)
		for i := 0; i < m; i++ {
			v, _ := strconv.Atoi(parts[1+i])
			arr[i] = v
		}
		pos := 1 + m
		q, _ := strconv.Atoi(parts[pos])
		pos++
		if len(parts) != pos+2*q {
			fmt.Printf("test %d: wrong number of query values\n", idx)
			os.Exit(1)
		}
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			l, _ := strconv.Atoi(parts[pos+i*2])
			r, _ := strconv.Atoi(parts[pos+i*2+1])
			queries[i] = [2]int{l, r}
		}
		// expected output
		exp := expected(arr, queries)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", m))
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", q))
		for i := 0; i < q; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", queries[i][0], queries[i][1]))
		}
		input := sb.String()

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(outFields) != len(exp) {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, len(exp), len(outFields))
			os.Exit(1)
		}
		for i, f := range outFields {
			val, _ := strconv.Atoi(f)
			if val != exp[i] {
				fmt.Printf("test %d failed\nexpected %v\ngot %v\n", idx, exp, outFields)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
