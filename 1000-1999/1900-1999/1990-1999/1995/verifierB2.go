package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func expectedB2(n, m int64, a, freq []int64) int64 {
	mp := make(map[int64]int64, n)
	for i := int64(0); i < n; i++ {
		mp[a[i]] = freq[i]
	}
	var mx int64
	for i := int64(0); i < n; i++ {
		var ans int64
		f := min64(m/a[i], mp[a[i]])
		f1 := min64((m-f*a[i])/(a[i]+1), mp[a[i]+1])
		ans = f*a[i] + f1*(a[i]+1)
		f3 := min64(f, mp[a[i]+1]-f1)
		ans += f3
		mx = max64(mx, min64(ans, m))
		if mx == m {
			break
		}
	}
	return mx
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesB2.txt")
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
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		var n int64
		var m int64
		fmt.Sscan(fields[0], &n)
		fmt.Sscan(fields[1], &m)
		if len(fields) != 2+int(2*n) {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		a := make([]int64, n)
		freq := make([]int64, n)
		for i := int64(0); i < n; i++ {
			fmt.Sscan(fields[2+i], &a[i])
		}
		for i := int64(0); i < n; i++ {
			fmt.Sscan(fields[2+int(n)+int(i)], &freq[i])
		}
		expect := expectedB2(n, m, a, freq)
		input := fmt.Sprintf("1\n%d %d\n", n, m)
		for i := int64(0); i < n; i++ {
			input += fmt.Sprintf("%d ", a[i])
		}
		input += "\n"
		for i := int64(0); i < n; i++ {
			input += fmt.Sprintf("%d ", freq[i])
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int64
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
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
