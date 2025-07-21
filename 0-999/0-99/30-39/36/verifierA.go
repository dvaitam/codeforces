package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, s string) string {
	pos := make([]int, 0, len(s))
	for i, ch := range s {
		if ch == '1' {
			pos = append(pos, i)
		}
	}
	if len(pos) < 3 {
		return "NO"
	}
	d := pos[1] - pos[0]
	for i := 2; i < len(pos); i++ {
		if pos[i]-pos[i-1] != d {
			return "NO"
		}
	}
	return "YES"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		var n int
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := scan.Text()
		expected := solveCase(n, s)

		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n%s\n", n, s)
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
