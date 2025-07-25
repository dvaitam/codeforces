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

func expected(n int, a, b []int) string {
	m := len(a)
	if m == 0 {
		return "YES"
	}
	u0, v0 := a[0], b[0]
	check := func(x int) bool {
		idx := make([]int, 0, m)
		for i := 0; i < m; i++ {
			if a[i] != x && b[i] != x {
				idx = append(idx, i)
			}
		}
		if len(idx) == 0 {
			return true
		}
		y1, y2 := a[idx[0]], b[idx[0]]
		ok := true
		for _, i := range idx {
			if a[i] != y1 && b[i] != y1 {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
		ok = true
		for _, i := range idx {
			if a[i] != y2 && b[i] != y2 {
				ok = false
				break
			}
		}
		return ok
	}
	if check(u0) || check(v0) {
		return "YES"
	}
	return "NO"
}

func runCase(exe string, input string, exp string) error {
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
		a := make([]int, m)
		b := make([]int, m)
		for i := 0; i < m; i++ {
			scan.Scan()
			a[i], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			b[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < m; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", a[i], b[i]))
		}
		in := sb.String()
		exp := expected(n, a, b) + "\n"
		if err := runCase(exe, in, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
