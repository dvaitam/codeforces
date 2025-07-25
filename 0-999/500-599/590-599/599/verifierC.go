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

func expected(arr []int) string {
	n := len(arr)
	prefixMax := make([]int, n)
	prefixMax[0] = arr[0]
	for i := 1; i < n; i++ {
		if arr[i] > prefixMax[i-1] {
			prefixMax[i] = arr[i]
		} else {
			prefixMax[i] = prefixMax[i-1]
		}
	}
	suffixMin := make([]int, n)
	suffixMin[n-1] = arr[n-1]
	for i := n - 2; i >= 0; i-- {
		if arr[i] < suffixMin[i+1] {
			suffixMin[i] = arr[i]
		} else {
			suffixMin[i] = suffixMin[i+1]
		}
	}
	ans := 1
	for i := 0; i < n-1; i++ {
		if prefixMax[i] <= suffixMin[i+1] {
			ans++
		}
	}
	return fmt.Sprintf("%d", ans)
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
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			arr[j], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[j]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
