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

func decStringMinusOne(s string) string {
	b := []byte(s)
	i := len(b) - 1
	for i >= 0 {
		if b[i] > '0' {
			b[i]--
			break
		}
		b[i] = '9'
		i--
	}
	j := 0
	for j < len(b)-1 && b[j] == '0' {
		j++
	}
	return string(b[j:])
}

func modPowInt(x int64, e int, m int64) int64 {
	res := int64(1)
	base := x % m
	for e > 0 {
		if e&1 != 0 {
			res = (res * base) % m
		}
		base = (base * base) % m
		e >>= 1
	}
	return res
}

func modPowDecimalExp(base int64, exp string, m int64) int64 {
	result := int64(1)
	for i := 0; i < len(exp); i++ {
		d := int(exp[i] - '0')
		result = modPowInt(result, 10, m)
		if d != 0 {
			result = (result * modPowInt(base, d, m)) % m
		}
	}
	return result
}

func solveCaseD(bStr, nStr string, c int64) int64 {
	var bMod int64
	for i := 0; i < len(bStr); i++ {
		bMod = (bMod*10 + int64(bStr[i]-'0')) % c
	}
	expStr := decStringMinusOne(nStr)
	p := modPowDecimalExp(bMod, expStr, c)
	bMinus1Mod := (bMod - 1 + c) % c
	totalMod := (bMinus1Mod * p) % c
	if totalMod == 0 {
		return c
	}
	return totalMod
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	type testCase struct {
		bStr string
		nStr string
		cStr string
		exp  string
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Printf("bad test file at case %d\n", i+1)
			os.Exit(1)
		}
		tests[i].bStr = scan.Text()
		if !scan.Scan() {
			fmt.Printf("bad test file at case %d\n", i+1)
			os.Exit(1)
		}
		tests[i].nStr = scan.Text()
		if !scan.Scan() {
			fmt.Printf("bad test file at case %d\n", i+1)
			os.Exit(1)
		}
		tests[i].cStr = scan.Text()
		cVal, _ := strconv.ParseInt(tests[i].cStr, 10, 64)
		tests[i].exp = fmt.Sprintf("%d", solveCaseD(tests[i].bStr, tests[i].nStr, cVal))
	}

	for i, tc := range tests {
		input := fmt.Sprintf("%s %s %s\n", tc.bStr, tc.nStr, tc.cStr)
		cmd := exec.Command(os.Args[1])
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("execution failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		outScan := bufio.NewScanner(bytes.NewReader(out))
		outScan.Split(bufio.ScanWords)
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != tc.exp {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected on test %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
