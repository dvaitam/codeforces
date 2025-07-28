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

func expected(n int, arr []int) string {
	sort.Ints(arr)
	m := len(arr)
	k := 0
	res := make([]string, n)
	for i := 1; i < n; i++ {
		res[i-1] = fmt.Sprintf("%d", arr[k])
		k += n - i
	}
	res[n-1] = fmt.Sprintf("%d", arr[m-1])
	return strings.Join(res, " ")
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
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("failed to read testcasesC.txt:", err)
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
			fmt.Printf("case %d missing array\n", caseNum)
			os.Exit(1)
		}
		fields := strings.Fields(scan.Text())
		arr := make([]int, len(fields))
		for i, f := range fields {
			arr[i], _ = strconv.Atoi(f)
		}
		input := fmt.Sprintf("1\n%d\n%s\n", n, strings.Join(fields, " "))
		want := expected(n, arr)
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
