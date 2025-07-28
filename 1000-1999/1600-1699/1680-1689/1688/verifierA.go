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

func expected(x int) int {
	if x == 1 {
		return 3
	}
	if x%2 == 1 {
		return 1
	}
	if x&(x-1) == 0 {
		return x + 1
	}
	return x & -x
}

func runCase(exe string, x int, exp int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(fmt.Sprintf("1\n%d\n", x))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotStr := strings.TrimSpace(out.String())
	if gotStr != fmt.Sprintf("%d", exp) {
		return fmt.Errorf("expected %d got %s", exp, gotStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
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
		x, _ := strconv.Atoi(scan.Text())
		exp := expected(x)
		if err := runCase(exe, x, exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
