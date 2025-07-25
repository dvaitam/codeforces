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

func solveCase(n int, l, r []int) string {
	sumL := 0
	sumR := 0
	for i := 0; i < n; i++ {
		sumL += l[i]
		sumR += r[i]
	}
	diff := abs(sumL - sumR)
	best := diff
	bestIdx := 0
	for i := 0; i < n; i++ {
		newL := sumL - l[i] + r[i]
		newR := sumR - r[i] + l[i]
		d := abs(newL - newR)
		if d > best {
			best = d
			bestIdx = i + 1
		}
	}
	return fmt.Sprintf("%d\n", bestIdx)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func runCase(exe string, input string, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
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
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		l := make([]int, n)
		r := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			l[j], _ = strconv.Atoi(scan.Text())
			scan.Scan()
			r[j], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", l[j], r[j]))
		}
		exp := solveCase(n, l, r)
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
