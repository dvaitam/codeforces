package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(a, b int64, xs []int64) []int64 {
	res := make([]int64, len(xs))
	for i, x := range xs {
		M := x * a / b
		w := (M*b + a - 1) / a
		res[i] = x - w
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		fmt.Println("bad test file:", err)
		os.Exit(1)
	}
	for caseID := 1; caseID <= t; caseID++ {
		var n int
		var a, b int64
		fmt.Fscan(reader, &n, &a, &b)
		xs := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &xs[i])
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
		for i, x := range xs {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprint(x))
		}
		input.WriteByte('\n')
		exp := expected(a, b, xs)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\n%s", caseID, err, out.String())
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out.String()))
		if len(fields) != n {
			fmt.Printf("test %d: expected %d numbers got %d\n", caseID, n, len(fields))
			os.Exit(1)
		}
		for i, f := range fields {
			var got int64
			fmt.Sscan(f, &got)
			if got != exp[i] {
				fmt.Printf("test %d failed: expected %v got %v\n", caseID, exp, fields)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
