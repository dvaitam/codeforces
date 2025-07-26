package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(x, y, z, a, b, c int) string {
	if a < x {
		return "NO"
	}
	if a+b < x+y {
		return "NO"
	}
	if a+b+c < x+y+z {
		return "NO"
	}
	return "YES"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("failed to read testcasesA.txt:", err)
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
		var x, y, z, a, b, c int
		fmt.Sscan(line, &x, &y, &z, &a, &b, &c)
		expect := strings.ToUpper(expected(x, y, z, a, b, c))
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
		got := strings.ToUpper(strings.TrimSpace(out.String()))
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
