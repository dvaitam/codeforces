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

func expected(xs []int64) string {
	n := len(xs)
	type pair struct {
		val int64
		idx int
	}
	arr := make([]pair, n)
	for i, v := range xs {
		arr[i] = pair{v, i}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].val < arr[j].val })
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] + arr[i].val
	}
	ans := make([]int64, n)
	total := int64(n)
	for i := 0; i < n; i++ {
		s := arr[i].val
		left := int64(i)*s - prefix[i]
		right := (prefix[n] - prefix[i+1]) - s*int64(n-i-1)
		ans[arr[i].idx] = total + left + right
	}
	out := make([]string, n)
	for i, v := range ans {
		out[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(out, " ")
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
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesE.txt")
	if err != nil {
		fmt.Println("failed to read testcasesE.txt:", err)
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
		if len(fields) != n {
			fmt.Printf("case %d wrong length\n", caseNum)
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i, f := range fields {
			arr[i], _ = strconv.ParseInt(f, 10, 64)
		}
		input := fmt.Sprintf("1\n%d\n%s\n", n, strings.Join(fields, " "))
		want := expected(arr)
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
