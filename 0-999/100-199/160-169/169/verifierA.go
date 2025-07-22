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

type testA struct {
	n, a, b int
	h       []int
}

func solveA(tc testA) int {
	h := append([]int(nil), tc.h...)
	sort.Ints(h)
	L := h[tc.b-1]
	R := h[tc.b]
	ans := R - L
	if ans < 0 {
		ans = 0
	}
	return ans
}

func parseLine(line string) (testA, error) {
	fields := strings.Fields(line)
	if len(fields) < 4 {
		return testA{}, fmt.Errorf("not enough fields")
	}
	n, _ := strconv.Atoi(fields[0])
	a, _ := strconv.Atoi(fields[1])
	b, _ := strconv.Atoi(fields[2])
	if len(fields) != 3+n {
		return testA{}, fmt.Errorf("expected %d numbers, got %d", n, len(fields)-3)
	}
	h := make([]int, n)
	for i := 0; i < n; i++ {
		h[i], _ = strconv.Atoi(fields[3+i])
	}
	return testA{n: n, a: a, b: b, h: h}, nil
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
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
		tc, err := parseLine(line)
		if err != nil {
			fmt.Printf("test %d parse error: %v\n", idx, err)
			os.Exit(1)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.a, tc.b))
		for i, v := range tc.h {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		expect := solveA(tc)
		got, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if got != fmt.Sprint(expect) {
			fmt.Printf("test %d failed: expected %d got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
