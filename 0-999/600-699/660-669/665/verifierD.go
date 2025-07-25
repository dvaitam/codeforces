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

func expected(arr []int) string {
	n := len(arr)
	maxVal := 0
	for _, v := range arr {
		if v > maxVal {
			maxVal = v
		}
	}
	limit := maxVal*2 + 2
	if limit < 3 {
		limit = 3
	}
	isPrime := make([]bool, limit)
	for i := 2; i < limit; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i < limit; i++ {
		if isPrime[i] {
			for j := i * i; j < limit; j += i {
				isPrime[j] = false
			}
		}
	}
	ans := make([]int, 0, n)
	for _, v := range arr {
		if v == 1 {
			ans = append(ans, 1)
		}
	}
	for _, v := range arr {
		if v > 1 && v+1 < limit && isPrime[v+1] {
			ans = append(ans, v)
			break
		}
	}
	if len(ans) > 1 {
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(len(ans)))
		sb.WriteByte('\n')
		for i, v := range ans {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		return sb.String()
	}
	ans = ans[:0]
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			s := arr[i] + arr[j]
			if s < limit && isPrime[s] {
				return fmt.Sprintf("2\n%d %d", arr[i], arr[j])
			}
		}
	}
	return fmt.Sprintf("1\n%d", arr[0])
}

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
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
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			arr[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		input := sb.String()
		exp := expected(arr) + "\n"
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseIdx+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
