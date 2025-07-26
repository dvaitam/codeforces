package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(k, a, b int64) string {
	if (a < k && b < k) || (a%k != 0 && b%k != 0) || (a%k != 0 && b < k) || (b%k != 0 && a < k) {
		return "-1"
	}
	return fmt.Sprintf("%d", a/k+b/k)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
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
		var k, a, b int64
		fmt.Sscan(scan.Text(), &k, &a, &b)
		expected := solveCase(k, a, b)
		input := fmt.Sprintf("%d %d %d\n", k, a, b)
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
