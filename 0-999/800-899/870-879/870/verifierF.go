package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func solveF(n int) int64 {
	spf := make([]int, n+1)
	phi := make([]int, n+1)
	primes := make([]int, 0, n/10)
	phi[1] = 1
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			spf[i] = i
			phi[i] = i - 1
			primes = append(primes, i)
		}
		for _, p := range primes {
			if p*i > n || p > spf[i] {
				break
			}
			spf[p*i] = p
			if i%p == 0 {
				phi[p*i] = phi[i] * p
				break
			} else {
				phi[p*i] = phi[i] * (p - 1)
			}
		}
	}
	var phiSum int64
	for i := 1; i <= n; i++ {
		phiSum += int64(phi[i])
	}
	totalPairs := int64(n*(n-1)) / 2
	edges := totalPairs - (phiSum - 1)

	limit := n / 2
	freq := map[int]int{}
	primesList := make([]int, 0)
	for i := 2; i <= n; i++ {
		p := spf[i]
		if p > limit {
			continue
		}
		if _, ok := freq[p]; !ok {
			primesList = append(primesList, p)
		}
		freq[p]++
	}
	sort.Ints(primesList)
	m := len(primesList)
	prefix := make([]int64, m+1)
	for i := m - 1; i >= 0; i-- {
		prefix[i] = prefix[i+1] + int64(freq[primesList[i]])
	}
	var count3 int64
	for i, p := range primesList {
		limitP := n / p
		j := sort.SearchInts(primesList, limitP+1)
		if j <= i {
			j = i + 1
		}
		if j < m {
			count3 += int64(freq[p]) * prefix[j]
		}
	}
	var sizeS int64
	for _, c := range freq {
		sizeS += int64(c)
	}
	pairS := sizeS * (sizeS - 1) / 2
	gcd1Pairs := pairS - edges
	count2 := gcd1Pairs - count3
	ans := edges + 2*count2 + 3*count3
	return ans
}

func runCase(exe, input, expect string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expect)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesF.txt")
	if err != nil {
		fmt.Println("could not read testcasesF.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		input := fmt.Sprintf("%d\n", n)
		expect := fmt.Sprintf("%d\n", solveF(n))
		if err := runCase(exe, input, expect); err != nil {
			fmt.Printf("case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
