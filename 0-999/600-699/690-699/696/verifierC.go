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

const MOD int64 = 1000000007
const PHI int64 = MOD - 1

func modPow(base, exp int64) int64 {
	res := int64(1)
	base %= MOD
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % MOD
		}
		base = base * base % MOD
		exp >>= 1
	}
	return res
}

func solveC(arr []int64) string {
	nMod := int64(1)
	parity := int64(1)
	for _, a := range arr {
		nMod = nMod * (a % PHI) % PHI
		parity = parity * (a % 2) % 2
	}
	exp := (nMod - 1 + PHI) % PHI
	pow2 := modPow(2, exp)
	inv3 := modPow(3, MOD-2)
	var numer int64
	if parity == 0 {
		numer = (pow2 + 1) % MOD
	} else {
		numer = (pow2 - 1 + MOD) % MOD
	}
	numer = numer * inv3 % MOD
	denom := pow2
	return fmt.Sprintf("%d/%d\n", numer, denom)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, k)
		var inSB strings.Builder
		inSB.WriteString(fmt.Sprintf("%d\n", k))
		for i := 0; i < k; i++ {
			scan.Scan()
			a, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[i] = a
			if i > 0 {
				inSB.WriteByte(' ')
			}
			inSB.WriteString(strconv.FormatInt(a, 10))
		}
		inSB.WriteByte('\n')
		input := inSB.String()
		expected := solveC(arr)
		if err := runCase(exe, input, expected); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
