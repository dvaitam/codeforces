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

type testCaseD struct {
	n   int
	arr []int64
}

func parseTestcasesD(path string) ([]testCaseD, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cases []testCaseD
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		vals := make([]int64, len(fields)-1)
		for i := range vals {
			v, _ := strconv.ParseInt(fields[1+i], 10, 64)
			vals[i] = v
		}
		cases = append(cases, testCaseD{n: n, arr: vals})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cases, nil
}

func solveD(n int, arr []int64) string {
	total := n + 2
	if len(arr) != total {
		return "-1"
	}
	var sum int64
	for _, v := range arr {
		sum += v
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	y := arr[total-1]
	x := sum - 2*y
	if idx := indexOf(arr, 0, total-1, x); idx >= 0 {
		var res []int64
		for i, v := range arr {
			if i == idx || i == total-1 {
				continue
			}
			res = append(res, v)
		}
		return joinInts(res)
	}
	y = arr[total-2]
	x = sum - 2*y
	if x == arr[total-1] {
		var res []int64
		for i, v := range arr {
			if i == total-2 || i == total-1 {
				continue
			}
			res = append(res, v)
		}
		return joinInts(res)
	}
	return "-1"
}

func joinInts(a []int64) string {
	if len(a) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	return sb.String()
}

func indexOf(a []int64, lo, hi int, x int64) int {
	idx := sort.Search(hi-lo, func(i int) bool { return a[lo+i] >= x })
	if idx < hi-lo && a[lo+idx] == x {
		return lo + idx
	}
	return -1
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestcasesD("testcasesD.txt")
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}
	for idx, tc := range cases {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		expected := solveD(tc.n, tc.arr)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
