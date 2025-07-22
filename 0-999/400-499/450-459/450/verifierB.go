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

func expectedValue(x, y, n int64) int64 {
	const mod int64 = 1000000007
	r := n % 6
	var res int64
	switch r {
	case 1:
		res = x
	case 2:
		res = y
	case 3:
		res = y - x
	case 4:
		res = -x
	case 5:
		res = -y
	case 0:
		res = x - y
	}
	res %= mod
	if res < 0 {
		res += mod
	}
	return res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
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
		if len(parts) != 3 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		x, _ := strconv.ParseInt(parts[0], 10, 64)
		y, _ := strconv.ParseInt(parts[1], 10, 64)
		nVal, _ := strconv.ParseInt(parts[2], 10, 64)
		expect := expectedValue(x, y, nVal)
		input := fmt.Sprintf("%d %d\n%d\n", x, y, nVal)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil {
			fmt.Printf("test %d: invalid output %s\n", idx, gotStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
