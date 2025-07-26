package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(arr []int) string {
	if len(arr) == 0 {
		return "0"
	}
	v := make([]int, 0, len(arr))
	prev := arr[0]
	v = append(v, prev)
	for i := 1; i < len(arr); i++ {
		if arr[i] != prev {
			v = append(v, arr[i])
			prev = arr[i]
		}
	}
	freq := make(map[int]int)
	maxF := 0
	for _, val := range v {
		freq[val]++
		if freq[val] > maxF {
			maxF = freq[val]
		}
	}
	ans := len(v) - maxF
	return fmt.Sprint(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to read testcasesD.txt:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(fields[0], &n)
		if len(fields) != 1+n {
			fmt.Printf("test %d expected %d numbers got %d\n", idx, 1+n, len(fields))
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(fields[1+i], &arr[i])
		}
		expect := solveCase(arr)
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
			fmt.Printf("test %d failed expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
