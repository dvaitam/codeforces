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

func expectedProfit(n, d int, costs []int, m int) int {
	sort.Ints(costs)
	used := m
	if used > n {
		used = n
	}
	sum := 0
	for i := 0; i < used; i++ {
		sum += costs[i]
	}
	if m > n {
		sum -= (m - n) * d
	}
	return sum
}

func runCase(bin string, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read testcasesA.txt: %v\n", err)
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
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		d, _ := strconv.Atoi(scan.Text())

		costs := make([]int, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
				os.Exit(1)
			}
			costs[i], _ = strconv.Atoi(scan.Text())
		}
		if !scan.Scan() {
			fmt.Fprintf(os.Stderr, "bad test case %d\n", caseNum)
			os.Exit(1)
		}
		mVal, _ := strconv.Atoi(scan.Text())

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, d)
		for i, v := range costs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		fmt.Fprintf(&sb, "%d\n", mVal)

		want := strconv.Itoa(expectedProfit(n, d, append([]int(nil), costs...), mVal))
		got, err := runCase(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", caseNum, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", caseNum, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", t)
}
