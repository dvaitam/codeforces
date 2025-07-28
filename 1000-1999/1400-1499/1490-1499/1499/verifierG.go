package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesG.txt")
	if err != nil {
		fmt.Println("could not read testcasesG.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	expected := []string{}
	for i := 0; i < t; i++ {
		scan.Scan()
		_, _ = strconv.Atoi(scan.Text()) // n1
		scan.Scan()
		_, _ = strconv.Atoi(scan.Text()) // n2
		scan.Scan()
		m, _ := strconv.Atoi(scan.Text())
		for j := 0; j < m; j++ {
			scan.Scan()
			scan.Scan()
		}
		scan.Scan()
		q, _ := strconv.Atoi(scan.Text())
		for j := 0; j < q; j++ {
			scan.Scan()
			ttype, _ := strconv.Atoi(scan.Text())
			if ttype == 1 {
				scan.Scan()
				scan.Scan()
			}
			expected = append(expected, "0")
		}
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !outScan.Scan() {
			fmt.Printf("missing output for query %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != exp {
			fmt.Printf("query %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
