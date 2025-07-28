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

func expected(n, m int, arr []int) int {
	cnt := make([]int, m)
	for _, v := range arr {
		cnt[v%m]++
	}
	ans := 0
	if cnt[0] > 0 {
		ans++
	}
	for r := 1; r*2 < m; r++ {
		a := cnt[r]
		b := cnt[m-r]
		if a == 0 && b == 0 {
			continue
		}
		if a > b {
			a, b = b, a
		}
		diff := b - a
		if diff <= 1 {
			ans++
		} else {
			ans += diff
		}
	}
	if m%2 == 0 {
		if cnt[m/2] > 0 {
			ans++
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	caseNum := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		caseNum++
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("case %d invalid format\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		m, _ := strconv.Atoi(fields[1])
		if len(fields)-2 != n {
			fmt.Printf("case %d invalid numbers count\n", caseNum)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				fmt.Printf("case %d invalid integer %q\n", caseNum, fields[2+i])
				os.Exit(1)
			}
			arr[i] = v
		}
		exp := expected(n, m, arr)
		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		resultStr := strings.TrimSpace(out.String())
		got, err := strconv.Atoi(resultStr)
		if err != nil {
			fmt.Printf("case %d: invalid output %q\n", caseNum, resultStr)
			os.Exit(1)
		}
		if got != exp {
			fmt.Printf("case %d failed: expected %d got %d\n", caseNum, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", caseNum)
}
