package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesB.txt")
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
		var a, b int
		fmt.Sscan(line, &a, &b)
		expected := a * b

		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(fmt.Sprintf("%d %d\n", a, b))
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		var got int
		fmt.Sscan(outStr, &got)
		if got != expected {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, expected, outStr)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", idx)
}
