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

func modInv(a int64) int64 {
	return modPow(a, MOD-2)
}

func comb(n, k int64) int64 {
	if k < 0 || k > n {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := int64(1); i <= k; i++ {
		res = res * ((n - k + i) % MOD) % MOD
		res = res * modInv(i%MOD) % MOD
	}
	return res
}

func expectedE(k int64, n, m int, A, B [][2]int64) int64 {
	points := []int64{0, k}
	for _, it := range A {
		points = append(points, it[0]-1, it[1])
	}
	for _, it := range B {
		points = append(points, it[0]-1, it[1])
	}
	sort.Slice(points, func(i, j int) bool { return points[i] < points[j] })
	uniq := []int64{}
	for _, p := range points {
		if len(uniq) == 0 || uniq[len(uniq)-1] != p {
			uniq = append(uniq, p)
		}
	}
	id := make(map[int64]int)
	for i, v := range uniq {
		id[v] = i
	}
	s := len(uniq)
	LA := make([][]int, s)
	LB := make([][]int, s)
	for _, it := range A {
		r := id[it[1]]
		l := id[it[0]-1]
		LA[r] = append(LA[r], l)
	}
	for _, it := range B {
		r := id[it[1]]
		l := id[it[0]-1]
		LB[r] = append(LB[r], l)
	}
	LBnd := make([]int64, s)
	UBnd := make([]int64, s)
	UBnd[0] = 0
	for i := 1; i < s; i++ {
		UBnd[i] = UBnd[i-1] + (uniq[i] - uniq[i-1])
	}
	for i := 1; i < s; i++ {
		if LBnd[i] < LBnd[i-1] {
			LBnd[i] = LBnd[i-1]
		}
		for _, l := range LA[i] {
			if LBnd[l]+1 > LBnd[i] {
				LBnd[i] = LBnd[l] + 1
			}
		}
		if UBnd[i] > UBnd[i-1]+(uniq[i]-uniq[i-1]) {
			UBnd[i] = UBnd[i-1] + (uniq[i] - uniq[i-1])
		}
		for _, l := range LB[i] {
			v := UBnd[l] + (uniq[i] - uniq[l] - 1)
			if v < UBnd[i] {
				UBnd[i] = v
			}
		}
		if LBnd[i] > UBnd[i] {
			return 0
		}
	}
	dp := map[int64]int64{0: 1}
	for i := 1; i < s; i++ {
		length := uniq[i] - uniq[i-1]
		newDP := map[int64]int64{}
		for sum, cnt := range dp {
			for add := int64(0); add <= length; add++ {
				ns := sum + add
				if ns < LBnd[i] || ns > UBnd[i] {
					continue
				}
				ways := comb(length, add)
				newDP[ns] = (newDP[ns] + cnt*ways) % MOD
			}
		}
		dp = newDP
	}
	ans := int64(0)
	for _, v := range dp {
		ans = (ans + v) % MOD
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		k, _ := strconv.ParseInt(parts[0], 10, 64)
		n, _ := strconv.Atoi(parts[1])
		m, _ := strconv.Atoi(parts[2])
		if len(parts) != 3+2*(n+m) {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, 3+2*(n+m), len(parts))
			os.Exit(1)
		}
		A := make([][2]int64, n)
		B := make([][2]int64, m)
		pos := 3
		for i := 0; i < n; i++ {
			l, _ := strconv.ParseInt(parts[pos], 10, 64)
			pos++
			r, _ := strconv.ParseInt(parts[pos], 10, 64)
			pos++
			A[i] = [2]int64{l, r}
		}
		for i := 0; i < m; i++ {
			l, _ := strconv.ParseInt(parts[pos], 10, 64)
			pos++
			r, _ := strconv.ParseInt(parts[pos], 10, 64)
			pos++
			B[i] = [2]int64{l, r}
		}
		expect := strconv.FormatInt(expectedE(k, n, m, A, B), 10)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
