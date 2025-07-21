package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

const modE = 51123987

func minE(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solveCaseE(s string) int64 {
	n := len(s)
	d1 := make([]int, n)
	l, r := 0, -1
	for i := 0; i < n; i++ {
		k := 1
		if i <= r {
			k = minE(d1[l+r-i], r-i+1)
		}
		for i-k >= 0 && i+k < n && s[i-k] == s[i+k] {
			k++
		}
		d1[i] = k
		if i+k-1 > r {
			l = i - k + 1
			r = i + k - 1
		}
	}
	d2 := make([]int, n)
	l, r = 0, -1
	for i := 0; i < n; i++ {
		k := 0
		if i <= r {
			k = minE(d2[l+r-i+1], r-i+1)
		}
		for i-k-1 >= 0 && i+k < n && s[i-k-1] == s[i+k] {
			k++
		}
		d2[i] = k
		if i+k-1 > r {
			l = i - k
			r = i + k - 1
		}
	}
	ds := make([]int64, n+2)
	de := make([]int64, n+2)
	var total int64
	for i := 0; i < n; i++ {
		if d1[i] > 0 {
			total += int64(d1[i])
			L := i - (d1[i] - 1)
			R := i
			ds[L]++
			ds[R+1]--
			de[i]++
			de[i+(d1[i]-1)+1]--
		}
		if d2[i] > 0 {
			total += int64(d2[i])
			L := i - d2[i]
			R := i - 1
			if L <= R {
				ds[L]++
				ds[R+1]--
			}
			de[i]++
			de[i+d2[i]]--
		}
	}
	cntStart := make([]int64, n)
	cntEnd := make([]int64, n)
	var cur int64
	for i := 0; i < n; i++ {
		cur += ds[i]
		cntStart[i] = cur
	}
	cur = 0
	for i := 0; i < n; i++ {
		cur += de[i]
		cntEnd[i] = cur
	}
	prefEnd := make([]int64, n)
	var cum int64
	for i := 0; i < n; i++ {
		cum = (cum + cntEnd[i]) % modE
		prefEnd[i] = cum
	}
	var disjoint int64
	for i := 0; i < n; i++ {
		cs := cntStart[i] % modE
		if cs != 0 && i > 0 {
			disjoint = (disjoint + cs*prefEnd[i-1]) % modE
		}
	}
	totalMod := total % modE
	inv2 := (modE + 1) / 2
	t := totalMod * ((totalMod - 1 + modE) % modE) % modE
	totalPairs := t * int64(inv2) % modE
	ans := (totalPairs - disjoint + modE) % modE
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
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
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		s := scan.Text()
		_ = n
		expected[i] = fmt.Sprintf("%d", solveCaseE(s))
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
