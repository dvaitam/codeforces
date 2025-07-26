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

func expected(arr []int) int {
	freq := make(map[int]int)
	maxCnt := 0
	for _, v := range arr {
		freq[v]++
		if freq[v] > maxCnt {
			maxCnt = freq[v]
		}
	}
	return maxCnt
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(arr)))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := expected(arr)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
