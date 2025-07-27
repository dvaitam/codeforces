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

func solveCase(a, b int64) int64 {
	r := (a + b) / 3
	if a < b {
		if r > a {
			r = a
		}
	} else {
		if r > b {
			r = b
		}
	}
	return r
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesA.txt")
	if err != nil {
		panic(err)
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
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("invalid test case %d\n", idx)
			os.Exit(1)
		}
		a, _ := strconv.ParseInt(fields[0], 10, 64)
		b, _ := strconv.ParseInt(fields[1], 10, 64)
		expected := solveCase(a, b)

		input := fmt.Sprintf("1\n%d %d\n", a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\n%s", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(out.String())
		if outStr != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
