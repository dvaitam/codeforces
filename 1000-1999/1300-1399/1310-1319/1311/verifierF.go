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

func expected(xs []int64, vs []int) string {
	n := len(xs)
	idx := make([]int, n)
	for i := 0; i < n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool { return xs[idx[i]] < xs[idx[j]] })
	var res int64
	for a := 0; a < n; a++ {
		for b := a + 1; b < n; b++ {
			i := idx[a]
			j := idx[b]
			if vs[i] <= vs[j] {
				res += xs[idx[b]] - xs[idx[a]]
			}
		}
	}
	return fmt.Sprintf("%d", res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idxT := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idxT++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d malformed\n", idxT)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+2*n {
			fmt.Printf("test %d malformed\n", idxT)
			os.Exit(1)
		}
		xs := make([]int64, n)
		vs := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+i])
			xs[i] = int64(v)
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+n+i])
			vs[i] = v
		}
		exp := expected(xs, vs)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", xs[i]))
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", vs[i]))
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
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idxT, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != exp {
			fmt.Printf("test %d failed\nexpected: %s\ngot: %s\n", idxT, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idxT)
}
