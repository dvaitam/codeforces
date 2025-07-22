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

const mod = 1000000007

func solve(n, h int, a []int) int64 {
	d := make([]int, n+2)
	for i := 1; i <= n; i++ {
		if a[i-1] > h {
			return 0
		}
		d[i] = h - a[i-1]
	}
	d[0], d[n+1] = 0, 0
	ans := int64(1)
	for i := 0; i <= n; i++ {
		cur := d[i]
		next := d[i+1]
		delta := next - cur
		switch {
		case delta == 1:
			// start
		case delta == 0:
			ans = ans * int64(cur+1) % mod
		case delta == -1:
			ans = ans * int64(cur) % mod
		default:
			return 0
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
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
		h, _ := strconv.Atoi(parts[1])
		if len(parts)-2 < n {
			fmt.Printf("test %d invalid len\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(parts[2+i])
			arr[i] = v
		}
		exp := solve(n, h, arr)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d %d\n", n, h)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, "%d ", arr[i])
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
