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

func expected(nums []int) string {
	var maxVal int
	var lastNeg int
	hasNeg := false
	for i, x := range nums {
		if x < 0 {
			hasNeg = true
			lastNeg = x
		}
		if i == 0 || x > maxVal {
			maxVal = x
		}
	}
	if hasNeg {
		return strconv.Itoa(lastNeg)
	}
	return strconv.Itoa(maxVal)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, "could not read testcasesA.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Fprintln(os.Stderr, "invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
				os.Exit(1)
			}
			nums[i], _ = strconv.Atoi(scan.Text())
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n", n)
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expectedOutput := expected(nums)
		var cmd *exec.Cmd
		if strings.HasSuffix(bin, ".go") {
			cmd = exec.Command("go", "run", bin)
		} else {
			cmd = exec.Command(bin)
		}
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", caseNum, err, errBuf.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expectedOutput {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, expectedOutput, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
