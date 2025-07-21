package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func colToLetters(c int) string {
	s := ""
	for c > 0 {
		c--
		s = string(rune('A'+(c%26))) + s
		c /= 26
	}
	return s
}

func lettersToCol(s string) int {
	c := 0
	for _, ch := range s {
		c = c*26 + int(ch-'A'+1)
	}
	return c
}

func isRXCY(s string) bool {
	if len(s) < 4 || s[0] != 'R' {
		return false
	}
	idx := strings.IndexRune(s, 'C')
	if idx == -1 {
		return false
	}
	for i := 1; i < idx; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	for i := idx + 1; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return idx > 1
}

func convert(s string) string {
	if isRXCY(s) {
		idx := strings.IndexRune(s, 'C')
		row := s[1:idx]
		col := s[idx+1:]
		c := 0
		fmt.Sscan(col, &c)
		return colToLetters(c) + row
	}
	// letters+row
	i := 0
	for i < len(s) && s[i] >= 'A' && s[i] <= 'Z' {
		i++
	}
	col := lettersToCol(s[:i])
	row := s[i:]
	return fmt.Sprintf("R%vC%v", row, col)
}

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
	var inputs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			inputs = append(inputs, line)
		}
	}
	expected := make([]string, len(inputs))
	for i, in := range inputs {
		expected[i] = convert(in)
	}
	inputStr := fmt.Sprintf("%d\n%s\n", len(inputs), strings.Join(inputs, "\n"))
	cmd := exec.Command(binary)
	cmd.Stdin = bytes.NewBufferString(inputStr)
	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	err = cmd.Run()
	if err != nil {
		fmt.Printf("runtime error: %v\nstderr: %s\n", err, errBuf.String())
		os.Exit(1)
	}
	outLines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	if len(outLines) != len(expected) {
		fmt.Printf("expected %d lines, got %d\n", len(expected), len(outLines))
		os.Exit(1)
	}
	for i, exp := range expected {
		if strings.TrimSpace(outLines[i]) != exp {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, exp, outLines[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(expected))
}
