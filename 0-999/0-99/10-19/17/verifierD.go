package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Base64-encoded contents of testcasesD.txt.
const testcasesD64 = "MTAwCjEwIDcgNzE0CjcgMTUgMTA2CjE0IDE4IDU4CjEwIDQgMTQxCjMgMSAyNjQKMTUgMTcgNjc3CjIwIDEwIDgyMgo0IDExIDQ5NAoyMCAxMSAyNzIKMTIgOCA0OQoxOCA4IDQ1MwoxOCAxIDY5NQoxNCAxMSA0NAoxOCA3IDczMwo1IDE4IDI1MAoyIDE4IDk5NgoxMiA1IDgxOAo5IDE5IDUwNQoxOSAxNiA0NDQKNyAxNSAyNjkKNyAyIDM3MAoxMSAxMSA2NDIKMTAgOCA0MDEKMiAyIDg3OQo1IDUgOTk4CjEyIDE5IDE2OQoyIDIgODE0CjE5IDYgOTY5CjIwIDUgODEKNiA1IDU2MQo3IDE0IDM5OQoyIDEyIDI5NgoxMyAxOSA3ODMKMiA4IDYwMwoyIDEyIDc5NAoxMyAxMSA4OQoyMCAxMiAzNzUKNCAxNCAxOTMKMTMgMTkgMjU5CjYgNyA4NTcKMTggNiA2MTUKNiAxNyA4MDMKNSAxMSA1OTYKMjAgMTggNjQxCjUgMjAgMTcwCjEzIDkgNDgwCjEwIDE1IDE5MAoxNCAxNCAzNzAKNSAxMSAxNTQKMTAgNiAzNDUKMTYgMyAzODQKMTMgOSA4OTEKMTMgNSA0MjQKOCA3IDE2MgoyMCAxNSA2MDkKMTUgMTkgNzIxCjkgMSAzNzUKNiAxIDQ1MgoxNSAxOSA4MjkKNSAyIDM5OAoxMSAyMCA3MDcKOCAxNCA0CjkgMTMgNjMwCjE3IDE2IDYxNQoxNSAxMyA5NTIKMTQgMiAzMTgKNyA2IDg4NAoxMiA0IDUxNwoxNiA2IDU5MAoxNCA4IDczNgo2IDkgMTkKMTYgMTEgNDkyCjEzIDE3IDc3NAo5IDEwIDU4MgoyIDE2IDcxMwo2IDkgMjEyCjEzIDEgODg0CjExIDIwIDQ0Mgo5IDEgMjI2CjE0IDUgODgxCjEwIDE5IDgyMgoxOSA3IDc5MAo2IDEzIDU3OQoxNyAxMiA5MTUKMTcgMjAgNzM2CjYgMSA3NjMKMTQgMTQgMzU0CjExIDUgMjk4CjE3IDkgNDA3CjkgMTUgODI4CjggNiA1MDkKMTYgMiA3NTUKMTUgNyA1NDAKMTEgMTQgMTA3CjExIDkgMTY3CjYgNyAzMTUKNiA2IDM1Ngo4IDE5IDYwOAo3IDIgMTA4CjYgOSA1ODkK"

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

type testCase struct {
	bStr string
	nStr string
	c    int64
}

func parseTestcases() ([]testCase, error) {
	raw, err := base64.StdEncoding.DecodeString(testcasesD64)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(raw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("invalid test data")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("parse t: %v", err)
	}
	tests := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		var tc testCase
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing b", i+1)
		}
		tc.bStr = scan.Text()
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing n", i+1)
		}
		tc.nStr = scan.Text()
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing c", i+1)
		}
		tc.c, err = strconv.ParseInt(scan.Text(), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("case %d parse c: %v", i+1, err)
		}
		tests = append(tests, tc)
	}
	if len(tests) != t {
		return nil, fmt.Errorf("expected %d cases, got %d", t, len(tests))
	}
	return tests, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	tests, err := parseTestcases()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for i, tc := range tests {
		exp := fmt.Sprintf("%d", solveCaseD(tc.bStr, tc.nStr, tc.c))
		input := fmt.Sprintf("%s %s %d\n", tc.bStr, tc.nStr, tc.c)
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
		if got != exp {
			fmt.Printf("test %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
		if outScan.Scan() {
			fmt.Printf("extra output detected on test %d\n", i+1)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
