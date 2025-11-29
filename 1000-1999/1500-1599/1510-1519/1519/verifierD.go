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

func solveD(a, b []int64) int64 {
	n := len(a)
	base := int64(0)
	for i := 0; i < n; i++ {
		base += a[i] * b[i]
	}
	ans := base
	for c := 0; c < n; c++ {
		l, r := c-1, c+1
		cur := base
		for l >= 0 && r < n {
			cur += a[l]*b[r] + a[r]*b[l] - a[l]*b[l] - a[r]*b[r]
			if cur > ans {
				ans = cur
			}
			l--
			r++
		}
	}
	for c := 0; c+1 < n; c++ {
		l, r := c, c+1
		cur := base
		for l >= 0 && r < n {
			cur += a[l]*b[r] + a[r]*b[l] - a[l]*b[l] - a[r]*b[r]
			if cur > ans {
				ans = cur
			}
			l--
			r++
		}
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
		fmt.Println(err)
		return
	}
	defer file.Close()
	in := bufio.NewReader(file)
	var t int
	fmt.Fscan(in, &t)
	for idx := 1; idx <= t; idx++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			a[i] = x
		}
		for i := 0; i < n; i++ {
			var x int64
			fmt.Fscan(in, &x)
			b[i] = x
		}
		expect := solveD(a, b)
		var input strings.Builder
		input.WriteString("")
		fmt.Fprintf(&input, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", a[i])
		}
		input.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", b[i])
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
			fmt.Printf("case %d runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != expect {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, expect, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
