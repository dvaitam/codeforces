package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	tokNum = iota
	tokComma
	tokDots
)

func expectedLine(s string) string {
	s = strings.TrimRight(s, "\r\n")
	var tokTypes []int
	var tokVals []string
	n := len(s)
	for i := 0; i < n; {
		switch {
		case s[i] >= '0' && s[i] <= '9':
			j := i
			for j < n && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			tokTypes = append(tokTypes, tokNum)
			tokVals = append(tokVals, s[i:j])
			i = j
		case i+2 < n && s[i] == '.' && s[i+1] == '.' && s[i+2] == '.':
			tokTypes = append(tokTypes, tokDots)
			tokVals = append(tokVals, "...")
			i += 3
		case s[i] == ',':
			tokTypes = append(tokTypes, tokComma)
			tokVals = append(tokVals, ",")
			i++
		case s[i] == ' ':
			i++
		default:
			i++
		}
	}
	var b strings.Builder
	m := len(tokTypes)
	for i, t := range tokTypes {
		switch t {
		case tokNum:
			if i > 0 && tokTypes[i-1] == tokNum {
				b.WriteByte(' ')
			}
			b.WriteString(tokVals[i])
		case tokComma:
			b.WriteString(",")
			if i < m-1 {
				b.WriteByte(' ')
			}
		case tokDots:
			if b.Len() > 0 {
				str := b.String()
				if str[len(str)-1] != ' ' {
					b.WriteByte(' ')
				}
			}
			b.WriteString("...")
		}
	}
	return b.String()
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
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		idx++
		expect := expectedLine(line)
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
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed\nexpected: %q\n     got: %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
