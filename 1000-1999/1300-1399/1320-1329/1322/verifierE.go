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

type testCase struct {
	n   int
	arr []int
}

func median3(a, b, c int) int {
	if a > b {
		if b > c {
			return b
		} else if a > c {
			return c
		}
		return a
	}
	if a > c {
		return a
	} else if b > c {
		return c
	}
	return b
}

func expected(tc testCase) (int, []int) {
	a := make([]int, tc.n)
	copy(a, tc.arr)
	steps := 0
	for {
		b := make([]int, tc.n)
		b[0] = a[0]
		if tc.n > 1 {
			b[tc.n-1] = a[tc.n-1]
		}
		for i := 1; i < tc.n-1; i++ {
			b[i] = median3(a[i-1], a[i], a[i+1])
		}
		same := true
		for i := 0; i < tc.n; i++ {
			if b[i] != a[i] {
				same = false
				break
			}
		}
		if same {
			break
		}
		a = b
		steps++
		if steps > 1000 {
			break
		}
	}
	return steps, a
}

func run(bin string, tc testCase) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", tc.n))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		parts := strings.Fields(line)
		n, _ := strconv.Atoi(parts[0])
		if len(parts)-1 != n {
			fmt.Fprintf(os.Stderr, "test %d malformed\n", idx)
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(parts[1+i])
		}
		tc := testCase{n: n, arr: arr}
		wantSteps, wantArr := expected(tc)
		gotStr, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		fields := strings.Fields(gotStr)
		if len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "case %d bad output\n", idx)
			os.Exit(1)
		}
		gotSteps, _ := strconv.Atoi(fields[0])
		if gotSteps != wantSteps {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, wantSteps, gotSteps)
			os.Exit(1)
		}
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+i])
			if v != wantArr[i] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %v got %v\n", idx, wantArr, fields[1:])
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
