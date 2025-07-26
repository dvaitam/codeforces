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

func expected(n, s int, times []int) string {
	if times[0] >= s+1 {
		return fmt.Sprintf("0 0\n")
	}
	for i := 0; i < n-1; i++ {
		if times[i+1]-times[i] >= 2*s+2 {
			t := times[i] + s + 1
			return fmt.Sprintf("%d %d\n", t/60, t%60)
		}
	}
	t := times[n-1] + s + 1
	return fmt.Sprintf("%d %d\n", t/60, t%60)
}

func runCase(bin, input, exp string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(exp), got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		return
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		return
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			return
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s, _ := strconv.Atoi(scan.Text())
		times := make([]int, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			h, _ := strconv.Atoi(scan.Text())
			scan.Scan()
			m, _ := strconv.Atoi(scan.Text())
			times[i] = h*60 + m
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, s))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", times[i]/60, times[i]%60))
		}
		in := sb.String()
		exp := expected(n, s, times)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, in)
			return
		}
	}
	fmt.Println("All tests passed")
}
