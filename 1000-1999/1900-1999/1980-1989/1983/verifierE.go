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

func modPow(a, b int64) int64 {
	res := int64(1)
	for b > 0 {
		if b&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		b >>= 1
	}
	return res
}

func expected(n, k int, vals []int64) string {
	sum := int64(0)
	for _, v := range vals {
		sum += v
	}
	N := n - k
	invN1 := modPow(int64(N+1), MOD-2)
	pSpec := (int64(N/2) + 1) % MOD * invN1 % MOD
	var pNorm int64
	if N > 0 {
		invN := modPow(int64(N), MOD-2)
		pNorm = (int64((N+1)/2) % MOD) * invN % MOD
	}
	alice := int64(0)
	for i, v := range vals {
		if i < k {
			alice = (alice + (v%MOD)*pSpec) % MOD
		} else {
			alice = (alice + (v%MOD)*pNorm) % MOD
		}
	}
	bob := (sum%MOD - alice) % MOD
	if bob < 0 {
		bob += MOD
	}
	if alice < 0 {
		alice += MOD
	}
	return fmt.Sprintf("%d %d", alice, bob)
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Println("failed to open testcasesE.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid testcases file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("missing n for case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		if !scan.Scan() {
			fmt.Printf("missing k for case %d\n", caseNum)
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		vals := make([]int64, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			vals[i] = v
		}
		expect := expected(n, k, vals)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", vals[i]))
		}
		sb.WriteByte('\n')
		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", caseNum, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
