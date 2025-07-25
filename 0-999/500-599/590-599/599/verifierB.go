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

func expected(n int, f []int, b []int) string {
	pos := make([][]int, n+1)
	for i, v := range f {
		pos[v] = append(pos[v], i+1)
	}
	ans := make([]int, len(b))
	ambiguous := false
	for i, v := range b {
		if len(pos[v]) == 0 {
			return "Impossible"
		}
		if len(pos[v]) > 1 {
			ambiguous = true
		}
		ans[i] = pos[v][0]
	}
	if ambiguous {
		return "Ambiguity"
	}
	var sb strings.Builder
	sb.WriteString("Possible\n")
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		fArr := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			fArr[i], _ = strconv.Atoi(scan.Text())
		}
		bArr := make([]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			bArr[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(fArr[i]))
		}
		sb.WriteByte('\n')
		for i := 0; i < m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(bArr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(n, fArr, bArr)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
