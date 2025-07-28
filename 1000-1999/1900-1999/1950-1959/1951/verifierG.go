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

const MOD int64 = 1_000_000_007

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= MOD
	for e > 0 {
		if e&1 == 1 {
			res = res * a % MOD
		}
		a = a * a % MOD
		e >>= 1
	}
	return res
}

func modInv(a int64) int64 { return modPow(a, MOD-2) }

func solve(n int, m int64, pos []int64) int64 {
	sort.Slice(pos, func(i, j int) bool { return pos[i] < pos[j] })
	var prefix, prefix2 int64
	var sum1, sum2 int64
	for j, v := range pos {
		x := v % MOD
		jm := int64(j) % MOD
		sum1 = (sum1 + (x*jm%MOD-MOD+prefix)%MOD) % MOD
		t := (x*x%MOD*jm%MOD - 2*x%MOD*prefix%MOD + prefix2) % MOD
		if t < 0 {
			t += MOD
		}
		sum2 = (sum2 + t) % MOD
		prefix = (prefix + x) % MOD
		prefix2 = (prefix2 + x*x%MOD) % MOD
	}
	total := (2 * ((m%MOD)*sum1%MOD - sum2 + MOD) % MOD) % MOD
	ans := total * modInv(int64(n)%MOD) % MOD
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesG.txt")
	if err != nil {
		panic(err)
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
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		mVal, _ := strconv.ParseInt(fields[1], 10, 64)
		pos := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[2+i], 10, 64)
			pos[i] = v
		}
		exp := solve(n, mVal, pos)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, mVal))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", pos[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != fmt.Sprintf("%d", exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
