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

func solveB(n, m int, a, b []string) string {
	words := make(map[string]bool, n)
	for _, w := range a {
		words[w] = true
	}
	common := 0
	for _, w := range b {
		if words[w] {
			common++
		}
	}
	if common%2 == 1 {
		n++
	}
	if n > m {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
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
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad test %d\n", idx)
			os.Exit(1)
		}
		n := atoi(fields[0])
		m := atoi(fields[1])
		if len(fields) != 2+n+m {
			fmt.Printf("bad test %d\n", idx)
			os.Exit(1)
		}
		a := fields[2 : 2+n]
		b := fields[2+n : 2+n+m]
		input := fmt.Sprintf("%d %d\n%s\n%s\n", n, m, strings.Join(a, " "), strings.Join(b, " "))
		expected := solveB(n, m, a, b)

		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}
