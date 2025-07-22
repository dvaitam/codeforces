package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(a []int) int {
	n := len(a)
	cnt := 0
	pos := -1
	for i := 0; i < n-1; i++ {
		if a[i] > a[i+1] {
			cnt++
			pos = i
		}
	}
	if cnt == 0 {
		return 0
	}
	if cnt > 1 || a[n-1] > a[0] {
		return -1
	}
	return n - pos - 1
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		var n int
		fmt.Sscan(fields[0], &n)
		if n <= 0 || len(fields)-1 < n {
			fmt.Printf("Test %d invalid test case\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[i+1], &arr[i])
		}
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:n+1], " "))
		exp := expected(arr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("Test %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(strings.TrimSpace(out), &got)
		if got != exp {
			fmt.Printf("Test %d wrong answer: expected %d got %s\n", idx, exp, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
