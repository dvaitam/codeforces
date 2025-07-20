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

func expected(n int, k int64, a []int64) int64 {
	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		b[i] = b[i-1] + a[i-1]
	}
	f := make([]int64, n+2)
	var ans int64
	j := n
	for i := n; i >= 1; i-- {
		for j > i && b[j-1]-b[i-1] > k {
			j--
		}
		if b[j]-b[i-1] > k {
			f[i] = int64(j-i) + f[j+1]
		} else {
			f[i] = int64(n - i + 1)
		}
		ans += f[i]
	}
	return ans
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		kval, _ := strconv.ParseInt(parts[1], 10, 64)
		if len(parts) != n+2 {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[2+i], 10, 64)
			arr[i] = v
		}
		expect := expected(n, kval, arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, kval))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("test %d: cannot parse output %q\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
