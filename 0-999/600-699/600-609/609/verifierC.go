package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expectedMoves(tasks []int64) int64 {
	n := len(tasks)
	sum := int64(0)
	for _, v := range tasks {
		sum += v
	}
	q := sum / int64(n)
	r := int(sum % int64(n))
	sort.Slice(tasks, func(i, j int) bool { return tasks[i] > tasks[j] })
	var diff int64
	for i := 0; i < n; i++ {
		target := q
		if i < r {
			target = q + 1
		}
		if tasks[i] > target {
			diff += tasks[i] - target
		} else {
			diff += target - tasks[i]
		}
	}
	return diff / 2
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 1 {
			fmt.Printf("test %d invalid\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		if len(parts)-1 < n {
			fmt.Printf("test %d invalid length\n", idx)
			os.Exit(1)
		}
		tasks := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(parts[1+i], 10, 64)
			tasks[i] = v
		}
		exp := expectedMoves(tasks)
		var buf bytes.Buffer
		fmt.Fprintf(&buf, "%d", n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&buf, " %d", tasks[i])
		}
		buf.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(buf.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != fmt.Sprint(exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
