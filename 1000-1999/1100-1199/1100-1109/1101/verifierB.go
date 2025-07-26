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

func solveCase(s string) int {
	n := len(s)
	st1 := -1
	for i := 0; i < n; i++ {
		if s[i] == '[' {
			st1 = i
			break
		}
	}
	if st1 < 0 {
		return -1
	}
	st2 := -1
	for i := st1 + 1; i < n; i++ {
		if s[i] == ':' {
			st2 = i
			break
		}
	}
	if st2 < 0 {
		return -1
	}
	en1 := -1
	for i := n - 1; i >= 0; i-- {
		if s[i] == ']' {
			en1 = i
			break
		}
	}
	if en1 < 0 || en1 <= st2 {
		return -1
	}
	en2 := -1
	for i := en1 - 1; i >= 0; i-- {
		if s[i] == ':' {
			en2 = i
			break
		}
	}
	if en2 < 0 || en2 <= st2 {
		return -1
	}
	cnt := 4
	for i := st2 + 1; i < en2; i++ {
		if s[i] == '|' {
			cnt++
		}
	}
	return cnt
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Scan()
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		s := scan.Text()
		expected := solveCase(strings.TrimSpace(s))
		input := s + "\n"
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
