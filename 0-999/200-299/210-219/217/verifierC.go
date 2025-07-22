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

func expected(line string) string {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ""
	}
	formula := parts[1]
	// Evaluate formula with ?=0 and ?=1
	op := make([]byte, 0, len(formula))
	val0 := make([]bool, 0, len(formula))
	val1 := make([]bool, 0, len(formula))
	for i := 0; i < len(formula); i++ {
		c := formula[i]
		switch c {
		case '0':
			val0 = append(val0, false)
			val1 = append(val1, false)
		case '1':
			val0 = append(val0, true)
			val1 = append(val1, true)
		case '?':
			val0 = append(val0, false)
			val1 = append(val1, true)
		case '(', '|', '&', '^':
			op = append(op, c)
		case ')':
			l := len(val0)
			v2_0 := val0[l-1]
			v2_1 := val1[l-1]
			val0 = val0[:l-1]
			val1 = val1[:l-1]
			l = len(val0)
			v1_0 := val0[l-1]
			v1_1 := val1[l-1]
			val0 = val0[:l-1]
			val1 = val1[:l-1]
			oper := op[len(op)-1]
			op = op[:len(op)-1]
			if len(op) > 0 && op[len(op)-1] == '(' {
				op = op[:len(op)-1]
			}
			var r0, r1 bool
			switch oper {
			case '|':
				r0 = v1_0 || v2_0
				r1 = v1_1 || v2_1
			case '&':
				r0 = v1_0 && v2_0
				r1 = v1_1 && v2_1
			case '^':
				r0 = v1_0 != v2_0
				r1 = v1_1 != v2_1
			}
			val0 = append(val0, r0)
			val1 = append(val1, r1)
		}
	}
	res0, res1 := false, false
	if len(val0) > 0 {
		res0 = val0[len(val0)-1]
		res1 = val1[len(val1)-1]
	}
	if res0 != res1 {
		return "YES"
	}
	return "NO"
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
		expect := expected(line)
		input := line + "\n"
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
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
