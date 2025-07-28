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

func expected(a, b []int64) string {
	n := len(a)
	diff := make([]int64, n)
	maxd := int64(-1 << 63)
	for i := 0; i < n; i++ {
		diff[i] = a[i] - b[i]
		if diff[i] > maxd {
			maxd = diff[i]
		}
	}
	var idx []string
	for i := 0; i < n; i++ {
		if diff[i] == maxd {
			idx = append(idx, fmt.Sprintf("%d", i+1))
		}
	}
	return fmt.Sprintf("%d\n%s", len(idx), strings.Join(idx, " "))
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
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to read testcasesD.txt:", err)
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
			fmt.Printf("case %d missing n\n", caseNum)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
		if !scan.Scan() {
			fmt.Printf("case %d missing a\n", caseNum)
			os.Exit(1)
		}
		fieldsA := strings.Fields(scan.Text())
		if len(fieldsA) != n {
			fmt.Printf("case %d wrong a length\n", caseNum)
			os.Exit(1)
		}
		arrA := make([]int64, n)
		for i, f := range fieldsA {
			arrA[i], _ = strconv.ParseInt(f, 10, 64)
		}
		if !scan.Scan() {
			fmt.Printf("case %d missing b\n", caseNum)
			os.Exit(1)
		}
		fieldsB := strings.Fields(scan.Text())
		if len(fieldsB) != n {
			fmt.Printf("case %d wrong b length\n", caseNum)
			os.Exit(1)
		}
		arrB := make([]int64, n)
		for i, f := range fieldsB {
			arrB[i], _ = strconv.ParseInt(f, 10, 64)
		}
		input := fmt.Sprintf("1\n%d\n%s\n%s\n", n, strings.Join(fieldsA, " "), strings.Join(fieldsB, " "))
		want := expected(arrA, arrB)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
