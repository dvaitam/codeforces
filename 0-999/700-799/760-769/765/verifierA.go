package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n int, home string, flights []string) string {
	diff := 0
	for _, f := range flights {
		if strings.HasPrefix(f, home) {
			diff++
		} else {
			diff--
		}
	}
	if diff == 0 {
		return "home"
	}
	return "contest"
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
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for caseNum := 1; caseNum <= t; caseNum++ {
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
		home := strings.TrimSpace(scan.Text())
		flights := make([]string, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			flights[i] = strings.TrimSpace(scan.Text())
		}
		expected := solveCase(n, home, flights)

		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n%s\n", n, home)
		for _, fl := range flights {
			input.WriteString(fl)
			input.WriteByte('\n')
		}

		cmd := exec.Command(os.Args[1])
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
