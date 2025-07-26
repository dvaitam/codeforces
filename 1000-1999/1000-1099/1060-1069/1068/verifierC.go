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

func solveCase(n int, edges [][2]int) string {
	sol := make([][]int, n+1)
	cnt := 0
	for _, e := range edges {
		cnt++
		sol[e[0]] = append(sol[e[0]], cnt)
		sol[e[1]] = append(sol[e[1]], cnt)
	}
	for i := 1; i <= n; i++ {
		if len(sol[i]) == 0 {
			cnt++
			sol[i] = append(sol[i], cnt)
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d\n", len(sol[i])))
		for _, id := range sol[i] {
			sb.WriteString(fmt.Sprintf("%d %d\n", i, id))
		}
	}
	return strings.TrimSpace(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line1 := strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}
		idx++
		parts := strings.Fields(line1)
		if len(parts) != 2 {
			fmt.Printf("invalid header on test %d\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		mVal, _ := strconv.Atoi(parts[1])
		if !scanner.Scan() {
			fmt.Printf("missing edge line for test %d\n", idx)
			os.Exit(1)
		}
		line2 := strings.TrimSpace(scanner.Text())
		nums := []string{}
		if line2 != "" {
			nums = strings.Fields(line2)
		}
		if len(nums) != 2*mVal {
			fmt.Printf("wrong number of edge values in test %d\n", idx)
			os.Exit(1)
		}
		edges := make([][2]int, mVal)
		for j := 0; j < mVal; j++ {
			x, _ := strconv.Atoi(nums[2*j])
			y, _ := strconv.Atoi(nums[2*j+1])
			edges[j] = [2]int{x, y}
		}
		// build input
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", nVal, mVal))
		for _, e := range edges {
			input.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		expected := solveCase(nVal, edges)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
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
		if got != expected {
			fmt.Printf("test %d failed\nexpected:\n%s\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
