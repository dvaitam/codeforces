package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, f []int) string {
	for i := 1; i <= n; i++ {
		if f[f[i]] != f[i] {
			return "-1"
		}
	}
	p := make(map[int]int)
	g := make([]int, n+1)
	h := make([]int, 0, n)
	s := 0
	for i := 1; i <= n; i++ {
		fi := f[i]
		idx, ok := p[fi]
		if !ok {
			s++
			idx = s
			p[fi] = s
			h = append(h, fi)
		}
		g[i] = idx
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", s)
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", g[i])
	}
	b.WriteByte('\n')
	for i := 0; i < s; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", h[i])
	}
	return b.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		fields := strings.Fields(scan.Text())
		if len(fields) != n {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		f := make([]int, n+1)
		for i, v := range fields {
			fmt.Sscan(v, &f[i+1])
		}
		expected := solveCase(n, f)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		for i := 1; i <= n; i++ {
			if i > 1 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", f[i])
		}
		input.WriteByte('\n')

		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
