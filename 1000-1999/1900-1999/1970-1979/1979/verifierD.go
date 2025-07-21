package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func solveCase(n, k int, s string) int {
	pre := make([]int, n+1)
	for i := 0; i < n; i++ {
		pre[i+1] = pre[i] + int(s[i]-'0')
	}
	f := make([]bool, n+1)
	f[n] = true
	for i := n - 1; i >= 0; i-- {
		j := i + k
		if j > n {
			j = n
		}
		ones := pre[j] - pre[i]
		allSame := ones == 0 || ones == j-i
		okNext := f[j]
		okChange := j == n || s[i] != s[j]
		f[i] = okNext && allSame && okChange
	}
	b0 := int(s[0] - '0')
	lastParity := ((n - 1) / k) % 2
	expectedLast := byte('0' + byte(b0^lastParity))
	ans := -1
	for p := 1; p <= n; p++ {
		parity := ((p - 1) / k) % 2
		expected := byte('0' + byte(b0^parity))
		if s[p-1] != expected {
			break
		}
		if f[p] && (p == n || s[p] == expectedLast) {
			ans = p
			break
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
	expected := make([]int, t)
	strs := make([]string, t)
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		expected[i] = solveCase(n, k, s)
		strs[i] = s
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
		got, err := strconv.Atoi(outScan.Text())
		if err != nil {
			fmt.Printf("bad output for test %d\n", i+1)
			os.Exit(1)
		}
		if got != expected[i] {
			fmt.Printf("test %d failed: expected %d got %d (string=%s)\n", i+1, expected[i], got, strs[i])
			os.Exit(1)
		}
	}
	if outScan.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed!")
}
