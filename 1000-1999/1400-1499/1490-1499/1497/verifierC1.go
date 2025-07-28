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

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int64) int64 {
	g := gcd(a, b)
	if g == 0 {
		return 0
	}
	res := a / g * b
	return res
}

func checkSolution(n int, parts []int64) bool {
	if len(parts) != 3 {
		return false
	}
	sum := int64(0)
	for _, v := range parts {
		if v <= 0 {
			return false
		}
		sum += v
	}
	if sum != int64(n) {
		return false
	}
	cur := parts[0]
	for i := 1; i < 3; i++ {
		cur = lcm(cur, parts[i])
		if cur > int64(n)/2 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC1.txt")
	if err != nil {
		fmt.Println("failed to open testcasesC1.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		caseNum++
		fields := strings.Fields(line)
		if len(fields) != 2 {
			fmt.Printf("case %d invalid format\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		if k != 3 {
			fmt.Printf("case %d invalid k\n", caseNum)
			os.Exit(1)
		}
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		outFields := strings.Fields(strings.TrimSpace(out.String()))
		if len(outFields) != 3 {
			fmt.Printf("case %d: expected 3 numbers got %d\n", caseNum, len(outFields))
			os.Exit(1)
		}
		parts := make([]int64, 3)
		for i, s := range outFields {
			v, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				fmt.Printf("case %d: invalid integer %q\n", caseNum, s)
				os.Exit(1)
			}
			parts[i] = v
		}
		if !checkSolution(n, parts) {
			fmt.Printf("case %d failed validation\n", caseNum)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}
