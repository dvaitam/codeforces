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

func hasSolution(n int, s, t string) bool {
	a0, a1, a2 := 0, 0, 0
	sumB := 0
	v := make([]int, n)
	for i := 0; i < n; i++ {
		a := int(s[i] - '0')
		b := int(t[i] - '0')
		if b == 1 {
			sumB++
		}
		val := a + b
		v[i] = val
		switch val {
		case 0:
			a0++
		case 1:
			a1++
		case 2:
			a2++
		}
	}
	half := n / 2
	for i2 := 0; i2 <= a2; i2++ {
		for j1 := 0; j1 <= a1; j1++ {
			if i2+j1 > half {
				continue
			}
			if 2*i2+j1 == sumB && half-i2-j1 <= a0 {
				return true
			}
		}
	}
	return false
}

func isValidOutput(n int, s, t string, out string) bool {
	out = strings.TrimSpace(out)
	if out == "-1" {
		return false
	}
	fields := strings.Fields(out)
	if len(fields) != n/2 {
		return false
	}
	seen := make(map[int]bool)
	firstClown := 0
	secondAcrobat := 0
	for _, f := range fields {
		idx, err := strconv.Atoi(f)
		if err != nil || idx < 1 || idx > n || seen[idx] {
			return false
		}
		seen[idx] = true
		if s[idx-1] == '1' {
			firstClown++
		}
	}
	for i := 0; i < n; i++ {
		if !seen[i+1] && t[i] == '1' {
			secondAcrobat++
		}
	}
	return firstClown == secondAcrobat
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
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n int
		var aStr, bStr string
		fmt.Sscan(line, &n, &aStr, &bStr)
		expectExists := hasSolution(n, aStr, bStr)
		input := fmt.Sprintf("%d\n%s\n%s\n", n, aStr, bStr)
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
		res := strings.TrimSpace(out.String())
		if expectExists {
			if res == "-1" || !isValidOutput(n, aStr, bStr, res) {
				fmt.Printf("test %d failed: invalid answer\n", idx)
				os.Exit(1)
			}
		} else {
			if res != "-1" {
				fmt.Printf("test %d failed: expected -1 got %s\n", idx, res)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
