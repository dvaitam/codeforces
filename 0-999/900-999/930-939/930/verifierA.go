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

func expectedA(n int, parents []int) int {
	depth := make([]int, n+1)
	count := map[int]int{}
	count[0] = 1
	for i := 2; i <= n; i++ {
		p := parents[i-2]
		depth[i] = depth[p] + 1
		count[depth[i]]++
	}
	ans := 0
	for _, c := range count {
		if c%2 == 1 {
			ans++
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts) != n {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, n, len(parts))
			os.Exit(1)
		}
		parents := make([]int, n-1)
		for i := 2; i <= n; i++ {
			v, _ := strconv.Atoi(parts[i-1])
			parents[i-2] = v
		}
		expect := strconv.Itoa(expectedA(n, parents))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
