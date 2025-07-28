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

func expected(x string) string {
	digits := make([]int, len(x)+1)
	for i := 0; i < len(x); i++ {
		digits[i+1] = int(x[i] - '0')
	}
	n := len(x)
	for i := n; i > 0; i-- {
		if digits[i] >= 5 {
			digits[i] = 0
			j := i - 1
			digits[j]++
			for j > 0 && digits[j] == 10 {
				digits[j] = 0
				j--
				digits[j]++
			}
			for k := i + 1; k <= n; k++ {
				digits[k] = 0
			}
		}
	}
	start := 0
	if digits[0] == 0 {
		start = 1
	}
	var b strings.Builder
	for i := start; i <= n; i++ {
		b.WriteByte(byte(digits[i]) + '0')
	}
	return b.String()
}

func runCase(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("failed to read testcasesB.txt:", err)
		os.Exit(1)
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseNum := 1; caseNum <= t; caseNum++ {
		if !scan.Scan() {
			fmt.Printf("case %d missing number\n", caseNum)
			os.Exit(1)
		}
		num := strings.TrimSpace(scan.Text())
		input := fmt.Sprintf("1\n%s\n", num)
		want := expected(num)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
