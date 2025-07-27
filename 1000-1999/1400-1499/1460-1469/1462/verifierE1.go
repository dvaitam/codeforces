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

func expected(arr []int) int64 {
	n := len(arr)
	sort.Ints(arr)
	var res int64
	l := 0
	for r := 0; r < n; r++ {
		for l < r && arr[r]-arr[l] > 2 {
			l++
		}
		cnt := r - l
		if cnt >= 2 {
			res += int64(cnt) * int64(cnt-1) / 2
		}
	}
	return res
}

func runCase(bin string, arr []int) error {
	n := len(arr)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintln(n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(arr[i]))
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.ParseInt(gotStr, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	exp := expected(arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE1.txt")
	if err != nil {
		fmt.Println("could not open testcasesE1.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scanner.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scanner.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scanner.Text())
			arr[j] = v
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
