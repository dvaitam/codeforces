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

func solve(n int, k int64, p []int) (string, []int) {
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[p[i-1]] = i
	}
	t := make([]int, n+2)
	add := func(x int) {
		for i := x; i <= n; i += i & -i {
			t[i]++
		}
	}
	ask := func(x int) int {
		s := 0
		for i := x; i > 0; i -= i & -i {
			s += t[i]
		}
		return s
	}
	tot := int64(n * (n - 1) / 2)
	sum := tot
	for val := 1; val <= n; val++ {
		sum -= int64(ask(a[val] - 1))
		add(a[val])
	}
	k -= sum
	if k%2 != 0 || k/2 < 0 || k/2 > (tot-sum) {
		return "NO", nil
	}
	k /= 2
	for i := range t {
		t[i] = 0
	}
	bArr := make([]int, n+1)
	for iVal := 1; iVal <= n; iVal++ {
		xCnt := ask(a[iVal] - 1)
		if k > int64(xCnt) {
			k -= int64(xCnt)
			add(a[iVal])
			continue
		}
		if k > 0 && k <= int64(xCnt) {
			id := 0
			for j := 1; j <= iVal; j++ {
				if a[j] < a[iVal] && k > 0 {
					id = j
					k--
				}
			}
			for j := 1; j <= id; j++ {
				bArr[j] = iVal + 1 - j
			}
			for j := id + 1; j <= iVal-1; j++ {
				bArr[j] = iVal - j
			}
			bArr[iVal] = iVal - id
			add(a[iVal])
			continue
		}
		bArr[iVal] = iVal
		add(a[iVal])
	}
	res := make([]int, n)
	for i := 1; i <= n; i++ {
		res[i-1] = bArr[i]
	}
	return "YES", res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	file, err := os.Open("testcasesF.txt")
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
		kVal, _ := strconv.ParseInt(fields[1], 10, 64)
		perm := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			perm[i] = v
		}
		status, arr := solve(n, kVal, perm)
		var exp string
		if status == "NO" {
			exp = "NO"
		} else {
			exp = "YES\n" + strings.TrimSpace(strings.Join(func() []string {
				s := make([]string, len(arr))
				for i, v := range arr {
					s[i] = strconv.Itoa(v)
				}
				return s
			}(), " "))
		}

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, kVal))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(perm[i]))
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
		if outStr != exp {
			fmt.Printf("Test %d failed: expected %q got %q\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
