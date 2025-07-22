package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n int, buttons []int) []int {
	res := make([]int, n+1)
	for _, b := range buttons {
		for i := b; i <= n; i++ {
			if res[i] == 0 {
				res[i] = b
			}
		}
	}
	return res[1:]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		buttons := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &buttons[i])
		}
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, b := range buttons {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprint(b))
		}
		input.WriteByte('\n')
		exp := expected(n, buttons)
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
			var got int
			fmt.Sscan(f, &got)
			if got != exp[i] {
				fmt.Printf("test %d failed: expected %v got %v\n", caseID, exp, fields)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
