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

func solveCase(n int) string {
	var ans float64
	for i := 1; i <= n; i++ {
		ans += 1.0 / float64(i)
	}
	return fmt.Sprintf("%.12f", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
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
	numbers := make([]int, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		numbers[i] = n
	}
	binary := os.Args[1]
	for i, n := range numbers {
		input := fmt.Sprintf("%d\n", n)
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\nstderr: %s\n", i+1, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(outBuf.String())
		expected := solveCase(n)
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
