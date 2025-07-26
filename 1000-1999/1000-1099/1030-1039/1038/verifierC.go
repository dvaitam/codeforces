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

func expected(n int, a, b []int) string {
	sa := append([]int(nil), a...)
	sb := append([]int(nil), b...)
	sort.Slice(sa, func(i, j int) bool { return sa[i] > sa[j] })
	sort.Slice(sb, func(i, j int) bool { return sb[i] > sb[j] })
	var sumA, sumB int64
	i, j := 0, 0
	turn := 0
	for i < n || j < n {
		if turn == 0 {
			if i < n && (j >= n || sa[i] > sb[j]) {
				sumA += int64(sa[i])
				i++
			} else {
				j++
			}
			turn = 1
		} else {
			if j < n && (i >= n || sb[j] > sa[i]) {
				sumB += int64(sb[j])
				j++
			} else {
				i++
			}
			turn = 0
		}
	}
	return fmt.Sprintf("%d", sumA-sumB)
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
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
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
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			a[i], _ = strconv.Atoi(scan.Text())
		}
		for i := 0; i < n; i++ {
			scan.Scan()
			b[i], _ = strconv.Atoi(scan.Text())
		}
		var sbInput strings.Builder
		sbInput.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sbInput.WriteByte(' ')
			}
			sbInput.WriteString(strconv.Itoa(a[i]))
		}
		sbInput.WriteByte('\n')
		for i := 0; i < n; i++ {
			if i > 0 {
				sbInput.WriteByte(' ')
			}
			sbInput.WriteString(strconv.Itoa(b[i]))
		}
		sbInput.WriteByte('\n')
		exp := expected(n, a, b)
		if err := runCase(exe, sbInput.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
