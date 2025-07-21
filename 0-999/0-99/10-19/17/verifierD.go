package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		bStr := scan.Text()
		scan.Scan()
		nStr := scan.Text()
		scan.Scan()
		cVal, _ := strconv.ParseInt(scan.Text(), 10, 64)
		expected[i] = fmt.Sprintf("%d", solveCaseD(bStr, nStr, cVal))
	}
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewReader(data)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("execution failed:", err)
		os.Exit(1)
	}
	outScan := bufio.NewScanner(bytes.NewReader(out))
	outScan.Split(bufio.ScanWords)
	for i := 0; i < t; i++ {
		if !outScan.Scan() {
			fmt.Printf("missing output for test %d\n", i+1)
			os.Exit(1)
		}
		got := outScan.Text()
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, expected[i], got)
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
