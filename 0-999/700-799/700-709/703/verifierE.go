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

const testcasesERaw = `100
5 9 12 17 1 15 8
6 2 6 4 12 16 8 13
5 4 19 8 1 7 14
3 6 13 6 3
2 20 20 15
2 5 1 1
2 25 7 6
2 10 11 7
5 29 7 6 7 13 10
1 12 14
2 30 5 9
1 11 10
5 19 1 20 11 3 10
3 27 10 16 11
2 16 16 6
1 9 1
6 12 13 1 18 14 12 13
5 27 1 15 2 6 20
2 4 8 15
3 17 12 17 9
4 4 19 12 10 2
4 3 7 11 17 20
3 30 5 11 9
6 18 3 10 11 10 6 3
6 5 10 16 6 2 3 20
5 30 13 2 8 20 12
3 15 14 5 2
6 2 16 11 7 5 19 5
6 29 14 4 6 14 12 5
1 28 14
3 5 15 20 6
5 15 16 11 16 9 10
4 13 5 4 13 18
2 21 16 11
2 3 16 9
5 26 18 17 12 3 12
6 19 2 10 12 18 9 16
3 25 10 11 6
5 28 1 16 18 9 11
6 9 15 10 17 12 12 9
6 12 14 12 6 15 12 11
5 5 17 6 7 12 16
3 23 3 14 6
5 25 19 17 14 10 20
5 25 9 1 7 6 19
4 20 6 8 6 2
4 8 6 2 5 4
3 6 16 7 18
1 14 15
3 13 20 3 19
2 8 12 1
3 13 9 14 4
6 27 18 12 2 18 20 10
1 10 18
5 11 19 10 12 5 14
4 27 19 18 12 15
2 6 20 13
5 16 7 5 20 3 12
6 1 13 4 11 19 20 18
2 11 19 13
4 14 8 16 10 16
6 13 13 6 20 20 9 10
4 9 14 1 11 10
4 30 10 5 16 1
1 22 20
4 8 10 2 5 13
1 16 18
5 9 8 16 2 8 16
3 28 5 10 10
4 20 16 17 20 4
1 25 5
3 10 18 11 20
3 24 17 1 15
3 12 19 5 2
1 9 18
4 22 4 18 7 1
4 25 14 20 19 16
4 16 13 7 10 15
1 10 1
6 25 14 19 10 16 10 5
2 16 18 16
3 18 5 14 19
5 26 2 3 8 9 3
1 22 1
3 24 14 3 13
6 16 2 4 4 8 20 4
6 29 5 10 15 5 6 20
2 14 6 3
5 7 2 18 4 13 3
3 2 19 19 4
6 13 20 5 1 14 3 11
6 20 16 16 12 12 2 5
6 10 5 19 17 10 18 18
5 8 9 3 18 8 9
3 17 5 8 12
4 24 13 6 5 1
6 11 3 19 2 3 4 17
5 15 8 13 15 16 11
1 26 17
1 18 13
`

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
	sc := bufio.NewScanner(strings.NewReader(testcasesERaw))
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
