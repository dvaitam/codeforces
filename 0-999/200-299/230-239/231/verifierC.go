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

type TestCaseC struct {
	n   int
	k   int64
	arr []int64
	ans string
}

func solveCaseC(n int, k int64, arr []int64) string {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + arr[i]
	}
	l := 0
	bestCount := 1
	bestValue := arr[0]
	for r := 0; r < n; r++ {
		for l <= r && int64(r-l+1)*arr[r]-(pref[r+1]-pref[l]) > k {
			l++
		}
		cnt := r - l + 1
		if cnt > bestCount || (cnt == bestCount && arr[r] < bestValue) {
			bestCount = cnt
			bestValue = arr[r]
		}
	}
	return fmt.Sprintf("%d %d", bestCount, bestValue)
}

func readCasesC(path string) ([]TestCaseC, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scan := bufio.NewScanner(bytes.NewReader(data))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("bad file")
	}
	t, _ := strconv.Atoi(scan.Text())
	cases := make([]TestCaseC, t)
	for i := 0; i < t; i++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		scan.Scan()
		k64, _ := strconv.ParseInt(scan.Text(), 10, 64)
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			scan.Scan()
			v, _ := strconv.ParseInt(scan.Text(), 10, 64)
			arr[j] = v
		}
		ans := solveCaseC(n, k64, append([]int64(nil), arr...))
		cases[i] = TestCaseC{n, k64, arr, ans}
	}
	return cases, nil
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := readCasesC("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
		expected := tc.ans
		got, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
