package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCaseA struct {
	n, x, y int
	arr     []int
}

func parseTestsA(path string) ([]testCaseA, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return nil, fmt.Errorf("empty test file")
	}
	blocks := strings.Split(raw, "\n\n")
	cases := make([]testCaseA, 0, len(blocks))
	for _, b := range blocks {
		lines := strings.Split(strings.TrimSpace(b), "\n")
		if len(lines) < 2 {
			continue
		}
		header := strings.Fields(lines[0])
		if len(header) != 3 {
			return nil, fmt.Errorf("invalid header: %q", lines[0])
		}
		n, _ := strconv.Atoi(header[0])
		x, _ := strconv.Atoi(header[1])
		y, _ := strconv.Atoi(header[2])
		nums := strings.Fields(lines[1])
		arr := make([]int, len(nums))
		for i, s := range nums {
			v, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			arr[i] = v
		}
		cases = append(cases, testCaseA{n: n, x: x, y: y, arr: arr})
	}
	return cases, nil
}

func solveA(tc testCaseA) int {
	n := tc.n
	a := tc.arr
	for d := 0; d < n; d++ {
		ok := true
		for j := d - tc.x; j < d && ok; j++ {
			if j >= 0 && a[j] < a[d] {
				ok = false
			}
		}
		for j := d + 1; j <= d+tc.y && ok; j++ {
			if j < n && a[j] < a[d] {
				ok = false
			}
		}
		if ok {
			return d + 1
		}
	}
	return -1
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases, err := parseTestsA("testcasesA.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for i, tc := range cases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.x, tc.y))
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveA(tc)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\noutput:\n%s", i+1, err, got)
			os.Exit(1)
		}
		outVal, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Printf("case %d: unable to parse output %q\n", i+1, got)
			os.Exit(1)
		}
		if outVal != expect {
			fmt.Printf("case %d failed: expected %d got %d\n", i+1, expect, outVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
