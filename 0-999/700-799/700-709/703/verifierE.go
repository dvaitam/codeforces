package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func factorize(x int64) (pr []int64, cnt []int) {
	for p := int64(2); p*p <= x; p++ {
		if x%p == 0 {
			c := 0
			for x%p == 0 {
				x /= p
				c++
			}
			pr = append(pr, p)
			cnt = append(cnt, c)
		}
	}
	if x > 1 {
		pr = append(pr, x)
		cnt = append(cnt, 1)
	}
	return
}

func exponents(v int64, primes []int64) []int {
	res := make([]int, len(primes))
	for i, p := range primes {
		for v%p == 0 {
			res[i]++
			v /= p
		}
	}
	return res
}

func expectedE(n int, k int64, arr []int64) (int, int64) {
	primes, cnts := factorize(k)
	m := len(primes)
	exps := make([][]int, n)
	for i := 0; i < n; i++ {
		exps[i] = exponents(arr[i], primes)
	}
	bestLen := n + 1
	var bestSum int64 = 1<<63 - 1
	total := 1 << n
	for mask := 1; mask < total; mask++ {
		l := bits.OnesCount(uint(mask))
		if l > bestLen {
			continue
		}
		var sum int64
		need := make([]int, m)
		for i := 0; i < n; i++ {
			if mask>>i&1 == 1 {
				sum += arr[i]
				for j := 0; j < m; j++ {
					need[j] += exps[i][j]
				}
			}
		}
		ok := true
		for j := 0; j < m; j++ {
			if need[j] < cnts[j] {
				ok = false
				break
			}
		}
		if !ok {
			continue
		}
		if l < bestLen || (l == bestLen && sum < bestSum) {
			bestLen = l
			bestSum = sum
		}
	}
	if bestLen == n+1 {
		return -1, 0
	}
	return bestLen, bestSum
}

func verifyOutput(n int, k int64, arr []int64, out string, expectLen int, expectSum int64) error {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	m, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("bad first number")
	}
	if m == -1 {
		if expectLen != -1 {
			return fmt.Errorf("expected len %d but got -1", expectLen)
		}
		return nil
	}
	if m != expectLen {
		return fmt.Errorf("expected len %d got %d", expectLen, m)
	}
	if len(fields) != m+1 {
		return fmt.Errorf("expected %d indices got %d", m, len(fields)-1)
	}
	sum := int64(0)
	prodDiv := int64(1)
	for i := 0; i < m; i++ {
		idx, err := strconv.Atoi(fields[1+i])
		if err != nil || idx < 1 || idx > n {
			return fmt.Errorf("bad index")
		}
		sum += arr[idx-1]
		prodDiv *= arr[idx-1]
	}
	if sum != expectSum {
		return fmt.Errorf("expected sum %d got %d", expectSum, sum)
	}
	if prodDiv%k != 0 {
		return fmt.Errorf("product not divisible by k")
	}
	return nil
}

func runCase(exe string, n int, k int64, arr []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expectLen, expectSum := expectedE(n, k, arr)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if err := verifyOutput(n, k, arr, out.String(), expectLen, expectSum); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("could not read testcasesE.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for idx := 0; idx < t; idx++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		kval, _ := strconv.ParseInt(sc.Text(), 10, 64)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			sc.Scan()
			val, _ := strconv.ParseInt(sc.Text(), 10, 64)
			arr[i] = val
		}
		if err := runCase(exe, n, kval, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
